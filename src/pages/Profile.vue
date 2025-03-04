<template>
  <div class="p-6 max-w-lg mx-auto">
    <h1 class="text-2xl font-bold mb-4">User Profile</h1>

    <div v-if="loading" class="text-gray-500">Loading profile...</div>

    <form v-else class="space-y-4" @submit.prevent="saveProfile">
      <div>
        <label class="block text-sm font-medium">Name</label>
        <input v-model="profile.name" type="text" class="w-full p-2 border rounded-lg">
      </div>

      <div>
        <label class="block text-sm font-medium">Email</label>
        <input v-model="profile.email" type="email" class="w-full p-2 border rounded-lg" disabled>
      </div>

      <div>
        <label class="block text-sm font-medium">Company</label>
        <input v-model="profile.company" type="text" class="w-full p-2 border rounded-lg">
      </div>

      <div>
        <label class="block text-sm font-medium">HEC URL</label>
        <input v-model="profile.hec_url" type="text" class="w-full p-2 border rounded-lg">
      </div>

      <div>
        <label class="block text-sm font-medium">HEC Token</label>
        <input v-model="profile.hec_token" type="password" class="w-full p-2 border rounded-lg">
      </div>

      <button type="submit" class="bg-blue-500 text-white px-4 py-2 rounded">
        Save Changes
      </button>
    </form>
  </div>
</template>

<script>
import { ref, onMounted } from "vue";
import { auth, db } from "@/firebase";
import { doc, getDoc, setDoc, updateDoc } from "firebase/firestore";
import { onAuthStateChanged } from "firebase/auth";

export default {
  setup() {
    const profile = ref({});
    const loading = ref(true);

    const loadProfile = async (user) => {
      try {
        if (!user) {
          console.error("❌ No authenticated user found.");
          alert("Please log in.");
          loading.value = false;
          return;
        }

        console.log(`📌 Fetching profile for user: ${user.uid}`);
        const profileRef = doc(db, "profiles", user.uid);
        const profileSnap = await getDoc(profileRef);

        if (profileSnap.exists()) {
          console.log("✅ Profile data found:", profileSnap.data());
          profile.value = profileSnap.data();
        } else {
          console.warn("⚠ No profile record found. Creating default profile.");
          profile.value = {
            name: user.displayName || "",
            email: user.email || "",
            company: "",
            hec_url: "",
            hec_token: "",
            image: "",
          };

          await setDoc(profileRef, profile.value);
          console.log("✅ Default profile created.");
        }
      } catch (error) {
        console.error("🔥 Error loading profile:", error);
        alert(`Failed to load profile: ${error.message}`);
      } finally {
        loading.value = false;
      }
    };

    // ✅ Function to save profile changes
    const saveProfile = async () => {
      if (!auth.currentUser) {
        console.error("❌ No authenticated user found.");
        return;
      }

      try {
        console.log("📌 Saving profile...");
        const profileRef = doc(db, "profiles", auth.currentUser.uid);
        await updateDoc(profileRef, profile.value);
        console.log("✅ Profile saved successfully!");
        alert("Profile updated successfully!");
      } catch (error) {
        console.error("🔥 Error saving profile:", error);
        alert("Failed to save profile.");
      }
    };

    onMounted(() => {
      onAuthStateChanged(auth, (user) => {
        if (user) {
          loadProfile(user);
        } else {
          loading.value = false;
        }
      });
    });

    return { profile, loading, saveProfile };
  },
};
</script>