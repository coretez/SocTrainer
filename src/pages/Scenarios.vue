<template>
  <div class="container mx-auto p-6">
    <h1 class="text-2xl font-bold mb-4">Scenarios</h1>

    <!-- Upload Button -->
    <button 
      @click="showUploadModal = true"
      class="mb-4 px-4 py-2 bg-indigo-500 text-white rounded-md hover:bg-indigo-600"
    >
      Upload Scenario
    </button>

    <!-- Show loading message -->
    <p v-if="loading" class="text-gray-500">Loading scenarios...</p>

    <!-- No scenarios found -->
    <p v-else-if="scenarios.length === 0" class="text-red-500">
      No scenarios found. Upload a new one.
    </p>

    <!-- Scenarios List -->
    <div v-else class="space-y-4">
      <div 
        v-for="scenario in scenarios" 
        :key="scenario.id" 
        class="p-4 border rounded-lg shadow-lg bg-white"
      >
        <h2 class="text-lg font-semibold">{{ scenario.name }}</h2>
        <p class="text-gray-600">{{ scenario.description }}</p>

        <!-- Buttons -->
        <div class="mt-2 flex space-x-2">
          <button 
            @click="viewScenario(scenario.stream)" 
            class="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600"
          >
            View JSON
          </button>

          <a 
            :href="scenario.stream" 
            target="_blank" 
            download 
            class="px-4 py-2 bg-green-500 text-white rounded-md hover:bg-green-600"
          >
            Download JSON
          </a>
        </div>
      </div>
    </div>

    <!-- Upload Modal -->
    <div v-if="showUploadModal" class="fixed inset-0 flex items-center justify-center bg-gray-800 bg-opacity-50">
      <div class="bg-white p-6 rounded-lg shadow-lg w-96">
        <h2 class="text-lg font-semibold mb-4">Upload Scenario</h2>

        <input v-model="newScenario.name" type="text" placeholder="Scenario Name" class="w-full p-2 border rounded mb-2" />
        <input v-model="newScenario.description" type="text" placeholder="Description" class="w-full p-2 border rounded mb-2" />
        
        <input type="file" @change="handleFileUpload" class="w-full p-2 border rounded mb-4" accept=".json" />

        <div v-if="uploading" class="text-blue-500">Uploading...</div>

        <div class="flex justify-between mt-2">
          <button @click="uploadScenario" class="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600" :disabled="uploading">
            Upload
          </button>
          <button @click="showUploadModal = false" class="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600">
            Cancel
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from "vue";
import { getFirestore, collection, getDocs, addDoc } from "firebase/firestore";
import { getAuth, onAuthStateChanged } from "firebase/auth";
import { getStorage, ref as storageRef, uploadBytes, getDownloadURL } from "firebase/storage";

export default {
  name: "ManageScenarios",
  setup() {
    const scenarios = ref([]);
    const loading = ref(true);
    const showUploadModal = ref(false);
    const newScenario = ref({ name: "", description: "" });
    const selectedFile = ref(null);
    const uploading = ref(false);
    const user = ref(null);

    const db = getFirestore();
    const auth = getAuth();
    const storage = getStorage();

    const fetchScenarios = async () => {
      if (!user.value) {
        console.warn("🔴 User not authenticated. Skipping fetch.");
        loading.value = false;
        return;
      }

      try {
        const querySnapshot = await getDocs(collection(db, "scenarios"));
        scenarios.value = querySnapshot.docs.map(doc => ({
          id: doc.id,
          ...doc.data(),
        }));
        console.log("✅ Scenarios loaded:", scenarios.value);
      } catch (error) {
        console.error("🔥 Error fetching scenarios:", error);
      } finally {
        loading.value = false;
      }
    };

    // Handle file selection
    const handleFileUpload = (event) => {
      selectedFile.value = event.target.files[0];
    };

    // Upload Scenario JSON
    const uploadScenario = async () => {
      if (!newScenario.value.name || !newScenario.value.description || !selectedFile.value) {
        alert("⚠ Please fill in all fields and select a file.");
        return;
      }

      if (!user.value) {
        console.error("❌ No authenticated user.");
        return;
      }

      try {
        uploading.value = true;

        // Upload file to Firebase Storage
        const filePath = `scenarios/${selectedFile.value.name}`;
        const fileRef = storageRef(storage, filePath);
        await uploadBytes(fileRef, selectedFile.value);
        const fileURL = await getDownloadURL(fileRef);

        // Save scenario details in Firestore **only after successful upload**
        const scenarioData = {
          name: newScenario.value.name,
          description: newScenario.value.description,
          file_name: selectedFile.value.name,
          stream: fileURL,
          createdBy: user.value.uid,
        };

        const docRef = await addDoc(collection(db, "scenarios"), scenarioData);
        scenarios.value.push({ id: docRef.id, ...scenarioData });

        // Reset form & close modal
        newScenario.value = { name: "", description: "" };
        selectedFile.value = null;
        showUploadModal.value = false;
        uploading.value = false;

        console.log("✅ Scenario uploaded successfully.");
      } catch (error) {
        console.error("🔥 Error uploading scenario:", error);
        alert("Error uploading scenario.");
      } finally {
        uploading.value = false;
      }
    };

    // View JSON scenario file
    const viewScenario = async (fileUrl) => {
      window.open(fileUrl, "_blank");
    };

    // Ensure authentication before loading scenarios
    onMounted(() => {
      onAuthStateChanged(auth, (loggedInUser) => {
        if (loggedInUser) {
          console.log(`✅ User authenticated: ${loggedInUser.uid}`);
          user.value = loggedInUser;
          fetchScenarios();
        } else {
          console.warn("❌ No user logged in.");
          loading.value = false;
        }
      });
    });

    return { 
      scenarios, loading, showUploadModal, newScenario, 
      selectedFile, uploading, fetchScenarios, 
      handleFileUpload, uploadScenario, viewScenario 
    };
  },
};
</script>

<style scoped>
.container {
  max-width: 800px;
}
</style>