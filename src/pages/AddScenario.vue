<template>
    <div class="container mx-auto p-6">
      <h1 class="text-2xl font-bold mb-4">Create New Scenario</h1>
  
      <div v-if="loading" class="text-gray-500">Loading profile...</div>
  
      <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Profile Info -->
        <div class="bg-white shadow-md rounded-lg p-6">
          <h2 class="text-xl font-semibold">Profile Information</h2>
          <p><strong>Search Site Token:</strong> {{ profile.ScenarioSiteToken }}</p>
          <p><strong>Base Search URL:</strong> {{ profile.ScenarioSiteUrl }}</p>
        </div>
  
        <!-- Scenario Details -->
        <div class="bg-white shadow-md rounded-lg p-6">
          <h2 class="text-xl font-semibold">Scenario Details</h2>
  
          <label>Scenario Name</label>
          <input v-model="newScenario.name" type="text" class="input-style" />
  
          <label>Scenario Description</label>
          <input v-model="newScenario.description" type="text" class="input-style" />
  
          <label>File Name</label>
          <input v-model="filename" type="text" class="input-style" />
  
          <button @click="getData" class="btn btn-blue">Get Data</button>
          <button @click="uploadScenario" class="btn btn-green" :disabled="uploading">
            {{ uploading ? "Uploading..." : "Upload Scenario" }}
          </button>
        </div>
      </div>
  
      <!-- Search Parameters -->
      <div class="mt-4 bg-white shadow-md rounded-lg p-6">
        <h2 class="text-lg font-semibold">Search Data</h2>
  
        <label>Search String</label>
        <input v-model="searchData.searchString" type="text" class="input-style" />
  
        <label>Grid Account</label>
        <input v-model="searchData.gridAccount" type="text" class="input-style" />
  
        <label>Start Time (milliseconds)</label>
        <input v-model="searchData.startTime" type="number" class="input-style" />
  
        <label>End Time (milliseconds)</label>
        <input v-model="searchData.endTime" type="number" class="input-style" />
      </div>
  
      <!-- Find & Replace Component -->
      <FindReplace :apiData="apiData" @update-api-data="updateApiData" />
  
      <!-- JSON Display -->
      <div v-if="apiData" class="mt-4">
        <h3 class="text-lg font-semibold">API Data</h3>
        <pre class="bg-gray-100 rounded p-3 overflow-auto text-xs">
          {{ JSON.stringify(apiData, null, 2) }}
        </pre>
      </div>
    </div>
  </template>
  
  <script>
import { ref, onMounted } from "vue";
import { getFirestore, doc, getDoc, addDoc, collection } from "firebase/firestore";
import { getAuth, onAuthStateChanged } from "firebase/auth";
import FindReplace from "@/components/FindReplace.vue"; 
import axios from "axios";

//Renamed Firebase ref
import { getStorage, ref as storageRef, uploadBytes, getDownloadURL } from "firebase/storage"; 


  
  export default {
    name: "AddScenario",
    components: { FindReplace },
    setup() {
      const profile = ref({});
      const loading = ref(true);
      const apiData = ref(null);
      const newScenario = ref({ name: "", description: "" });
      const filename = ref("scenario.json");
      const searchData = ref({
        searchString: "",
        startTime: null,
        endTime: null,
        gridAccount: "",
        site: "",
        token: "",
      });
      const storage = getStorage();
      const uploading = ref(false);
      const user = ref(null);
  
      const db = getFirestore();
      const auth = getAuth();
  
      // 🔹 Fetch user profile and update searchData
      const fetchProfile = async () => {
        if (!user.value) {
          console.warn("🔴 User not authenticated.");
          loading.value = false;
          return;
        }
  
        try {
          const docRef = doc(db, "profiles", user.value.uid);
          const docSnap = await getDoc(docRef);
          if (docSnap.exists()) {
            profile.value = docSnap.data();
            searchData.value.site = profile.value.ScenarioSiteUrl;
            searchData.value.token = profile.value.ScenarioSiteToken;
          } else {
            profile.value = { message: "No profile found" };
          }
        } catch (error) {
          console.error("🔥 Error fetching profile:", error);
        } finally {
          loading.value = false;
        }
      };
  
      // 🔹 Fetch data from the backend
      const getData = async () => {
        try {
          console.log("🔍 Fetching data with:", searchData.value);
          const response = await axios.post("http://localhost:8080/api/get-data", searchData.value);
          apiData.value = response.data;
        } catch (error) {
          console.error("🔥 Error fetching data:", error);
        }
      };
  
      // 🔹 Update API Data (after FindReplace modifications)
      const updateApiData = (updatedData) => {
        apiData.value = updatedData;
      };
  
      const uploadScenario = async () => {
        if (!apiData.value || !newScenario.value.name || !newScenario.value.description) {
            alert("⚠ Please fill in all fields and ensure data is available.");
            return;
        }

        if (!user.value) {
            console.error("❌ No authenticated user.");
            return;
        }

        try {
            uploading.value = true;
            console.log("📤 Uploading scenario to Go backend...");

            // Convert apiData to a JSON Blob
            const jsonBlob = new Blob([JSON.stringify(apiData.value)], { type: "application/json" });

            // Create a FormData object
            const formData = new FormData();
            formData.append("file", jsonBlob, filename.value);
            formData.append("name", newScenario.value.name);
            formData.append("description", newScenario.value.description);
            formData.append("createdBy", user.value.uid);

            // Send request to Go backend
            const response = await axios.post("http://localhost:8080/api/upload-scenario", formData, {
                headers: { "Content-Type": "multipart/form-data" },
            });

            console.log("✅ Scenario uploaded:", response.data.file_url);
            alert(`Scenario uploaded successfully: ${response.data.file_url}`);

            // Reset form
            newScenario.value = { name: "", description: "" };
            apiData.value = null;
            uploading.value = false;
        } catch (error) {
            console.error("🔥 Error uploading scenario:", error);
            alert("Error uploading scenario.");
        } finally {
            uploading.value = false;
        }
    };
  
      onMounted(() => {
        onAuthStateChanged(auth, (loggedInUser) => {
          if (loggedInUser) {
            user.value = loggedInUser;
            fetchProfile();
          } else {
            console.warn("❌ No user logged in.");
            loading.value = false;
          }
        });
      });
  
      return {
        profile, loading, searchData, apiData, newScenario, filename,
        uploading, getData, uploadScenario, updateApiData
      };
    },
  };
  </script>
  
  <style scoped>
  .input-style { 
    @apply shadow border rounded w-full py-2 px-3 text-gray-700; 
  }
  .btn { 
    @apply py-2 px-4 rounded font-bold focus:outline-none; 
  }
  .btn-blue { 
    @apply bg-blue-500 hover:bg-blue-700 text-white; 
  }
  .btn-green { 
    @apply bg-green-500 hover:bg-green-700 text-white; 
  }
  </style>