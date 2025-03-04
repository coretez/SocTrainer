<template>
  <nav class="bg-gray-800 text-white p-4 flex justify-between">
    <h1 class="text-xl font-bold">
      Replay Trainer
    </h1>
    <ul class="flex space-x-4">
      <li>
        <router-link
          to="/"
          class="hover:underline"
        >
          Home
        </router-link>
      </li>
      <li>
        <router-link
          to="/scenarios"
          class="hover:underline"
        >
          Scenarios
        </router-link>
      </li>
      <li>
        <router-link
          to="/profile"
          class="hover:underline"
        >
          Profile
        </router-link>
      </li>
      <li>
        <router-link
          to="/workbook"
          class="hover:underline"
        >
          Workbook
        </router-link>
      </li>
      <li>
        <router-link
          v-if="!isAuthenticated"
          to="/login"
          class="hover:underline"
        >
          Login
        </router-link>
      </li>
      <li v-if="isAuthenticated">
        <button
          class="bg-red-500 px-3 py-1 rounded hover:bg-red-600"
          @click="logout"
        >
          Logout
        </button>
      </li>
    </ul>
  </nav>
</template>

<script>
import { getAuth, onAuthStateChanged, signOut } from 'firebase/auth';
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';

export default {
  setup() {
    const isAuthenticated = ref(false);
    const auth = getAuth();
    const router = useRouter();

    onMounted(() => {
      onAuthStateChanged(auth, (user) => {
        isAuthenticated.value = !!user;
      });
    });

    const logout = async () => {
      await signOut(auth);
      localStorage.removeItem("authenticated");
      isAuthenticated.value = false;
      router.push("/login");
    };

    return { isAuthenticated, logout };
  }
};
</script>