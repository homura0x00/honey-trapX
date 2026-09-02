<template>
  <div id="login-view" class="form-card2">
    <form class="form" @submit.prevent="handleSubmit">
      <p class="form-heading">Log In</p>

      <div class="form-field">
        <input
            required
            placeholder="email"
            class="input-field"
            type="email"
            v-model="formData.email"
        />
      </div>

      <div class="form-field">
        <input
            required
            placeholder="Password"
            class="input-field"
            type="password"
            v-model="formData.password"
        />
      </div>

      <div class="error-message" v-if="formData.errorMessage">
        {{ formData.errorMessage }}
      </div>

      <button class="sendMessage-btn" type="submit">Sign In</button>
    </form>
  </div>
</template>

<script setup lang="ts">
import {reactive} from "vue";

const formData = reactive({
  email: "",
  password: "",
  errorMessage: "",
});

const handleSubmit = () => {
  // 简单验证
  if (!formData.email || !formData.password) {
    formData.errorMessage = "请填写所有必填字段";
    return;
  }

  // 验证邮箱格式
  const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!emailPattern.test(formData.email)) {
    formData.errorMessage = "请输入有效的邮箱地址";
    return;
  }

  // 清除错误消息
  formData.errorMessage = "";

  // 在实际应用中，这里会发送API请求到后端
  console.log("提交登录表单: ", {
    email: formData.email,
    password: formData.password,
  });

  // 模拟登录成功
  alert("[TEST] 登录成功！");
};
</script>

<style scoped>
.form {
  display: flex;
  flex-direction: column;
  align-self: center;
  font-family: inherit;
  gap: 10px;
  padding-inline: 2em;
  padding-bottom: 0.4em;
  background-color: #171717;
  /* background-color: #0a192f; */
  border-radius: 20px;
}

.form-heading {
  text-align: center;
  margin: 2em;
  color: #64ffda;
  font-size: 1.2em;
  background-color: transparent;
  align-self: center;
}

.form-field {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5em;
  border-radius: 10px;
  padding: 0.6em;
  border: none;
  outline: none;
  color: white;
  background-color: #171717;
  box-shadow: inset 2px 5px 10px rgb(5, 5, 5);
}

.input-field {
  background: none;
  border: none;
  outline: none;
  width: 100%;
  color: #ccd6f6;
  padding-inline: 1em;
}

.sendMessage-btn {
  cursor: pointer;
  margin-bottom: 3em;
  padding: 1em;
  border-radius: 10px;
  border: none;
  outline: none;
  background-color: transparent;
  color: #64ffda;
  font-weight: bold;
  outline: 1px solid #64ffda;
  transition: all ease-in-out 0.3s;
}

.sendMessage-btn:hover {
  transition: all ease-in-out 0.3s;
  background-color: #64ffda;
  color: #000;
  cursor: pointer;
  box-shadow: inset 2px 5px 10px rgb(5, 5, 5);
}

.form-card2 {
  border-radius: 0;
  transition: all 0.2s;
}

.form-card2:hover {
  transform: scale(0.98);
  border-radius: 20px;
}

</style>
