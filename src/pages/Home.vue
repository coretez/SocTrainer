<template>
  <div class="p-6 max-w-4xl mx-auto">
    <h1 class="text-2xl font-bold mb-4">Welcome to Replay Trainer</h1>

    <!-- Loading Indicator -->
    <p v-if="loading" class="text-gray-500">Loading workbooks...</p>

    <!-- No Workbooks Message -->
    <p v-else-if="totalWorkbooks === 0" class="text-gray-500">
      No workbooks available yet.
    </p>

    <!-- Progress Summary -->
    <div v-else class="mb-6 p-4 bg-gray-100 rounded-md shadow">
      <h2 class="text-lg font-semibold text-indigo-700 mb-2">Your Progress</h2>
      <p class="text-gray-700">
        <span class="font-bold text-green-600">Completed:</span> {{ completedWorkbooks.length }}/{{ totalWorkbooks }}
      </p>
      <p class="text-gray-700">
        <span class="font-bold text-yellow-600">In Progress:</span> {{ workbooksInProgress.length }}/{{ totalWorkbooks }}
      </p>
      <p class="text-gray-700">
        <span class="font-bold text-red-600">Not Started:</span> {{ totalWorkbooks - (completedWorkbooks.length + workbooksInProgress.length) }}
      </p>
    </div>

    <!-- Workbooks Status -->
    <div v-if="totalWorkbooks > 0">
      <!-- In-Progress Workbooks -->
      <div v-if="workbooksInProgress.length > 0" class="mb-6">
        <h2 class="text-lg font-semibold text-yellow-600 mb-2">Workbooks In Progress</h2>
        <ul class="space-y-2">
          <li v-for="workbook in workbooksInProgress" :key="workbook.id" class="bg-yellow-100 p-3 rounded-md">
            <router-link :to="`/workbook/${workbook.id}`" class="text-blue-600 hover:underline">
              {{ workbook.scenarioName || "Unknown Scenario" }}
            </router-link>
            <p class="text-sm text-gray-600 prose max-w-none" v-html="workbook.scenarioDescription"></p>
          </li>
        </ul>
      </div>

      <!-- Completed Workbooks -->
      <div v-if="completedWorkbooks.length > 0">
        <h2 class="text-lg font-semibold text-green-600 mb-2">Completed Workbooks</h2>
        <ul class="space-y-2">
          <li v-for="workbook in completedWorkbooks" :key="workbook.id" class="bg-green-100 p-3 rounded-md">
            <router-link :to="`/workbook/${workbook.id}`" class="text-blue-600 hover:underline">
              {{ workbook.scenarioName || "Unknown Scenario" }}
            </router-link>
            <p class="text-sm text-gray-600 prose max-w-none" v-html="workbook.scenarioDescription"></p>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from "vue";
import { collection, getDocs } from "firebase/firestore";
import { auth, db } from "@/firebase";

export default {
  name: "Home",
  setup() {
    const workbooksInProgress = ref([]);
    const completedWorkbooks = ref([]);
    const totalWorkbooks = ref(0);
    const loading = ref(true);

    // Fetch user progress from Firestore
    const fetchUserProgress = async () => {
      const user = auth.currentUser;
      if (!user) {
        loading.value = false;
        return;
      }

      try {
        // Get total workbooks available
        const workbooksSnap = await getDocs(collection(db, "workbooks"));
        totalWorkbooks.value = workbooksSnap.size; // Total number of workbooks

        // Get user progress
        const progressRef = collection(db, "progress", user.uid, "workbooks");
        const progressSnap = await getDocs(progressRef);

        workbooksInProgress.value = [];
        completedWorkbooks.value = [];

        progressSnap.forEach((doc) => {
          const data = doc.data();
          if (data.status === "completed") {
            completedWorkbooks.value.push({ id: doc.id, ...data });
          } else {
            workbooksInProgress.value.push({ id: doc.id, ...data });
          }
        });

        console.log("📌 Workbooks In Progress:", workbooksInProgress.value);
        console.log("📌 Completed Workbooks:", completedWorkbooks.value);
        console.log("📌 Total Workbooks:", totalWorkbooks.value);
      } catch (error) {
        console.error("🔥 Error fetching progress:", error);
      } finally {
        loading.value = false;
      }
    };

    onMounted(fetchUserProgress);

    return { workbooksInProgress, completedWorkbooks, totalWorkbooks, loading };
  },
};
</script>