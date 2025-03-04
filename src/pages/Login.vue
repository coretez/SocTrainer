<template>
  <div class="bg-gray-50 dark:bg-gray-900">
    <div class="flex flex-col items-center justify-center px-6 py-8 mx-auto md:h-screen lg:py-0">
      <a
        href="#"
        class="flex items-center mb-6 text-2xl font-semibold text-gray-900 dark:text-white"
      >
        <img
          class="w-8 h-8 mr-2"
          src="https://tpc.googlesyndication.com/simgad/16390710324977171557"
          alt="logo"
        >
        Fluency    
      </a>
      <div class="w-full bg-white rounded-lg shadow dark:border md:mt-0 sm:max-w-md xl:p-0 dark:bg-gray-800 dark:border-gray-700">
        <div class="p-6 space-y-4 md:space-y-6 sm:p-8">
          <h1 class="text-xl font-bold leading-tight tracking-tight text-gray-900 md:text-2xl dark:text-white">
            Sign in to your account
          </h1>
          <form
            class="space-y-4 w-full"
            @submit.prevent="login"
          >
            <div>
              <label
                for="email"
                class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
              >Your email</label>
              <input
                id="username"
                v-model="username"
                type="email"
                class="w-full p-2 border rounded-lg"
                placeholder="name@company.com"
                required
              >
            </div>
            <div>
              <label
                for="password"
                class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
              >Password</label>
              <input
                id="password"
                v-model="password"
                type="password"
                class="w-full p-2 border rounded-lg"
                placeholder="••••••••"
                required
              >
            </div>
            <button
              type="submit"
              class="w-full bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700"
            >
              Sign in
            </button>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { login } from "@/services/authService";

export default {
  name: "LoginPage",
  data() {
    return {
      username: "",
      password: "",
    };
  },
  methods: {
    async login() {
      try {
        const user = await login(this.username, this.password);
        console.log("Logged in user:", user);
        localStorage.setItem("authenticated", "true");
        this.$router.push("/scenarios");
      } catch (error) {
        console.error("Error logging in:", error.message);
        alert("Failed to log in. Please check your email and password.");
      }
    },
  },
};
</script>
