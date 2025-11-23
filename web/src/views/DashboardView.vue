<template>
  <div class="dashboard">
    <header>
      <h2>My Prayers</h2>
      <button @click="handleLogout" class="logout-btn">Log Out</button>
    </header>

    <main>
      <div v-if="isLoading">Loading user profile...</div>
      <div v-else class="welcome-banner">
        <h3>Welcome, {{ user?.name }}</h3>
        <p>Your spiritual journey continues here.</p>
        <!-- Prayer list will go here later -->
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { useAuth0 } from '@auth0/auth0-vue';

const { logout, user, isLoading } = useAuth0();

const handleLogout = () => {
  logout({
    logoutParams: {
      returnTo: window.location.origin
    }
  });
};
</script>

<style scoped>
.dashboard {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
}

header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  border-bottom: 1px solid #e5e7eb;
  padding-bottom: 1rem;
}

.logout-btn {
  background: transparent;
  border: 1px solid #ef4444;
  color: #ef4444;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
}

.logout-btn:hover {
  background: #fef2f2;
}
</style>
