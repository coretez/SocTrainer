<template>
  <div class="p-6 max-w-4xl mx-auto">
    <h1 class="text-2xl font-bold mb-4">Workbook Details</h1>

    <div v-if="loading" class="text-gray-500">Loading workbook...</div>
    <div v-else-if="!workbook" class="text-red-500">Workbook not found.</div>

    <div v-else class="space-y-6">
      <!-- Scenario Title -->
      <h1 class="text-2xl font-bold mb-4">{{ workbook.scenarioName }}</h1>

      <p class="text-gray-700" v-html="workbook.fields.scenarioIntroduction"></p>

      <div class="mt-4">
        <label class="block text-sm font-bold">HEC URL</label>
        <input v-model="hecUrl" class="border rounded w-full p-2 mt-1" disabled />

        <label class="block text-sm font-bold mt-2">HEC Token</label>
        <input v-model="hecToken" type="password" class="border rounded w-full p-2 mt-1" disabled />
      </div>

      <button 
        @click="startReplay"
        class="mt-4 bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
        :disabled="replayInProgress"
      >
        {{ replayInProgress ? "Replaying..." : "Start Replay" }}
      </button>

      <!--p v-if="message" class="mt-2 text-green-600">{{ message }}</p -->
      <p v-if="errorMessage" class="mt-2 text-red-600">{{ errorMessage }}</p>

      <!-- Scenario Replay Progress -->
      <div v-if="progress.total > 0">
        <h3 class="text-lg font-semibold mt-4">Replay Progress</h3>
        <p>Record: {{ progress.rec }} / {{ progress.total }}</p>
        <p>Timestamp: {{ new Date(progress.timestamp * (progress.timestamp < 1e12 ? 1000 : 1)).toLocaleString() }}</p>

        <!-- Progress Bar -->
        <div class="w-full bg-gray-200 rounded h-4 mt-2">
          <div
            class="bg-blue-500 h-4 rounded"
            :style="{ width: (progress.rec / progress.total) * 100 + '%' }"
          ></div>
        </div>
      </div>

      <!-- Workbook Sections -->
      <div>
        <h3 class="text-lg font-semibold mb-1">Scenario Introduction</h3>
        <div class="prose max-w-none" v-html="workbook.fields.scenarioIntroduction"></div>
      </div>

      <div>
        <h3 class="text-lg font-semibold mb-1">Validation Stage</h3>
        <div class="prose max-w-none" v-html="workbook.fields.validationStage"></div>
      </div>

      <div>
        <h3 class="text-lg font-semibold mb-1">Scoping Stage</h3>
        <div class="prose max-w-none" v-html="workbook.fields.scopingStage"></div>
      </div>

      <div>
        <h3 class="text-lg font-semibold mb-1">Recovery & Response Stage</h3>
        <div class="prose max-w-none" v-html="workbook.fields.recoveryResponseStage"></div>
      </div>

      <div>
        <h3 class="text-lg font-semibold mb-1">Hotwash</h3>
        <div class="prose max-w-none" v-html="workbook.fields.hotwash"></div>
      </div>

      <!-- Navigation Buttons -->
      <div class="mt-4 flex space-x-2">
        <button
          v-if="!isStarted"
          @click="startWorkbook"
          class="px-4 py-2 bg-green-500 text-white rounded-md hover:bg-green-600"
        >
          Start Workbook
        </button>

        <router-link
          v-else
          :to="`/edit-workbook/${route.params.id}`"
          class="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600"
        >
          Continue Workbook
        </router-link>

        <router-link
          to="/workbook"
          class="px-4 py-2 bg-gray-500 text-white rounded-md hover:bg-gray-600"
        >
          Back to Workbooks
        </router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import { doc, getDoc, setDoc, serverTimestamp } from "firebase/firestore";
import { auth, db } from "@/firebase";
import axios from "axios";

export default {
  name: "WorkbookDetails",
  setup() {
    const route = useRoute();
    const workbook = ref(null);
    const hecUrl = ref("");
    const hecToken = ref("");
    const message = ref("");
    const errorMessage = ref("");
    const loading = ref(true);
    const isStarted = ref(false);
    const replayInProgress = ref(false);
    const progress = ref({ rec: 0, total: 0, timestamp: 0 });

    // Fetch Profile Data (HEC URL & Token)
    const fetchUserProfile = async () => {
      const user = auth.currentUser;
      if (!user) {
        console.error("User not logged in.");
        return;
      }

      try {
        const profileRef = doc(db, "profiles", user.uid);
        const profileSnap = await getDoc(profileRef);

        if (profileSnap.exists()) {
          const profileData = profileSnap.data();
          hecUrl.value = profileData.hec_url || "";
          hecToken.value = profileData.hec_token || "";
        } else {
          console.error("Profile not found.");
        }
      } catch (error) {
        console.error("Error fetching profile:", error);
      }
    };

    // Fetch Workbook Details
    const fetchWorkbook = async () => {
      try {
        const docRef = doc(db, "workbooks", route.params.id);
        const docSnap = await getDoc(docRef);

        if (docSnap.exists()) {
          workbook.value = docSnap.data();
          await checkWorkbookProgress();
        } else {
          console.error("Workbook not found.");
        }
      } catch (error) {
        console.error("Error fetching workbook:", error);
      } finally {
        loading.value = false;
      }
    };

    // Start Replay
    const startReplay = async () => {
      if (!hecUrl.value || !hecToken.value) {
        errorMessage.value = "HEC URL and Token are required!";
        return;
      }

      replayInProgress.value = true;
      message.value = "Replay started...";
      errorMessage.value = "";

      const requestData = {
        hec_token: hecToken.value,
        hec_url: hecUrl.value,
        scenario_name: workbook.value.scenarioName,
      };

      try {
        await axios.post("http://localhost:8080/api/replay", requestData);
        message.value = "Replay started successfully!";
        listenForProgress();
      } catch (error) {
        errorMessage.value = "Failed to start replay.";
        console.error("Replay error:", error);
        replayInProgress.value = false;
      }
    };

    // Listen for progress updates via SSE
    const listenForProgress = () => {
      const eventSource = new EventSource("http://localhost:8080/api/replay/progress");

      eventSource.onmessage = (event) => {
        const data = JSON.parse(event.data);
        progress.value = data;

        if (data.rec === data.total) {
          eventSource.close();
          replayInProgress.value = false;
        }
      };

      eventSource.onerror = (error) => {
        console.error("Error receiving progress updates.", error);
        eventSource.close();
        replayInProgress.value = false;
      };
    };

    // Check if workbook is started
    const checkWorkbookProgress = async () => {
      const user = auth.currentUser;
      if (!user) return;

      const progressRef = doc(db, "progress", user.uid, "workbooks", route.params.id);
      const progressSnap = await getDoc(progressRef);

      if (progressSnap.exists()) {
        isStarted.value = true;
      }
    };

    // Fetch Data on Mount
    onMounted(async () => {
      await fetchUserProfile();
      await fetchWorkbook();
    });

    return {
      route,
      workbook,
      loading,
      isStarted,
      hecUrl,
      hecToken,
      message,
      errorMessage,
      replayInProgress,
      progress,
      startReplay,
    };
  },
};
</script>