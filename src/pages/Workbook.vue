<template>
  <div class="p-6 max-w-4xl mx-auto">
    <h1 class="text-2xl font-bold mb-4">Workbooks</h1>

    <!-- Create Workbook Button -->
    <div class="mb-4">
      <router-link 
        to="/create-workbook" 
        class="px-4 py-2 bg-green-500 text-white rounded-md hover:bg-green-600"
      >
        + Create Workbook
      </router-link>
    </div>

    <!-- Show loading message -->
    <p v-if="loading" class="text-gray-500">Loading workbooks...</p>

    <!-- No workbooks found -->
    <p v-else-if="workbooks.length === 0" class="text-red-500">
      No workbooks found. Create one now.
    </p>

    <!-- Workbooks List -->
    <div v-else>
      <!-- In Progress Workbooks -->
      <div v-if="workbooksInProgress.length">
        <h2 class="text-xl font-semibold text-yellow-600">In Progress</h2>
        <div v-for="workbook in workbooksInProgress" :key="workbook.id" class="p-4 border rounded-lg shadow-lg bg-white mt-2">
          <h3 class="text-lg font-semibold">{{ workbook.scenarioName || "Unknown Scenario" }}</h3>
          <p class="text-gray-600">Status: {{ workbook.status }}</p>
          <br>
          <router-link :to="`/workbook/${workbook.id}`" class="mt-2 px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600">
            Continue
          </router-link>
        </div>
      </div>

      <!-- Completed Workbooks -->
      <div v-if="completedWorkbooks.length" class="mt-6">
        <h2 class="text-xl font-semibold text-green-600">Completed</h2>
        <div v-for="workbook in completedWorkbooks" :key="workbook.id" class="p-4 border rounded-lg shadow-lg bg-white mt-2">
          <h3 class="text-lg font-semibold">{{ workbook.scenarioName || "Unknown Scenario" }}</h3>
          <p class="text-gray-600">Completed on: {{ new Date(workbook.completedAt?.seconds * 1000).toLocaleDateString() }}</p>
          <br>
          <router-link :to="`/workbook/${workbook.id}`" class="mt-2 px-4 py-2 bg-gray-500 text-white rounded-md hover:bg-gray-600">
            View
          </router-link>
        </div>
      </div>

      <!-- Workbooks List (All Available Workbooks) -->
      <div v-if="workbooks.length" class="mt-6">
        <h2 class="text-xl font-semibold text-indigo-600">All Workbooks</h2>
        <div v-for="workbook in workbooks" :key="workbook.id" class="p-4 border rounded-lg shadow-lg bg-white mt-2">
          <h3 class="text-lg font-semibold">{{ workbook.scenarioName || "Unknown Scenario" }}</h3>
          <div class="text-gray-700 prose max-w-none" v-html="workbook.fields.scenarioIntroduction"></div>
          <br>
          <router-link :to="`/workbook/${workbook.id}`" class="mt-2 px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600">
            View Workbook
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

---

### **Updated Script**
- Fetches **all workbooks** from `workbooks/`
- Retrieves user progress from `progress/{user}/workbooks`
- Displays workbooks grouped into **In Progress**, **Completed**, and **All Workbooks**

```vue
<script>
import { ref, onMounted } from "vue";
import { collection, getDocs } from "firebase/firestore";
import { auth, db } from "@/firebase";

export default {
  name: "Workbooks",
  setup() {
    const workbooks = ref([]); // All workbooks
    const workbooksInProgress = ref([]);
    const completedWorkbooks = ref([]);
    const loading = ref(true);

    const fetchWorkbooks = async () => {
      try {
        const querySnapshot = await getDocs(collection(db, "workbooks"));
        workbooks.value = querySnapshot.docs.map(doc => ({
          id: doc.id,
          ...doc.data(),
        }));
      } catch (error) {
        console.error("Error fetching workbooks:", error);
      }
    };

    const fetchUserProgress = async () => {
      const user = auth.currentUser;
      if (!user) return;

      const progressRef = collection(db, "progress", user.uid, "workbooks");
      const progressSnap = await getDocs(progressRef);

      workbooksInProgress.value = [];
      completedWorkbooks.value = [];

      progressSnap.forEach((doc) => {
        const data = doc.data();
        if (data.status === "completed") {
          completedWorkbooks.value.push({
            id: doc.id,
            scenarioName: data.scenarioName || "Unknown Scenario",
            scenarioDescription: data.scenarioDescription || "",
            completedAt: data.completedAt,
          });
        } else {
          workbooksInProgress.value.push({
            id: doc.id,
            scenarioName: data.scenarioName || "Unknown Scenario",
            scenarioDescription: data.scenarioDescription || "",
            status: data.status,
          });
        }
      });

      console.log("📌 Workbooks in progress:", workbooksInProgress.value);
      console.log("📌 Completed workbooks:", completedWorkbooks.value);
    };

    onMounted(async () => {
      await fetchWorkbooks();
      await fetchUserProgress();
      loading.value = false;
    });

    return { workbooks, workbooksInProgress, completedWorkbooks, loading };
  },
};
</script>