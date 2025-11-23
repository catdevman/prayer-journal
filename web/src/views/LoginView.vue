<template>
  <div class="login-container">
    <div class="card">
      <h1>Prayer Journal</h1>
      <p>A secure place for your spiritual journey.</p>

      <div v-if="isLoading" class="loading">
        Loading...
      </div>

      <div v-else>
        <button v-if="!isAuthenticated" @click="handleLogin" class="btn-primary">
          Log In / Sign Up
        </button>
        <button v-else @click="goToDashboard" class="btn-secondary">
          Go to Dashboard
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuth0 } from '@auth0/auth0-vue';
import { useRouter } from 'vue-router';
import { watchEffect } from 'vue';

const { loginWithRedirect, isAuthenticated, isLoading } = useAuth0();
const router = useRouter();

const handleLogin = () => {
  loginWithRedirect({
    appState: { target: '/dashboard' }
  });
};

const goToDashboard = () => {
  router.push('/dashboard');
};

// Auto-redirect if they land here while already logged in
watchEffect(() => {
    if (!isLoading.value && isAuthenticated.value) {
        router.push('/dashboard');
    }
});
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f3f4f6;
}
.card {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  text-align: center;
}
.btn-primary {
  background-color: #2563eb;
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  border: none;
  cursor: pointer;
  font-size: 1rem;
}
.btn-primary:hover {
  background-color: #1d4ed8;
}
</style>
