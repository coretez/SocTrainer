<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold mb-4">
      Add Sample Workbook
    </h1>

    <form
      class="space-y-4"
      @submit.prevent="createWorkbook"
    >
      <div>
        <label class="block text-sm font-medium text-gray-700">Scenario ID</label>
        <input
          v-model="scenarioId"
          type="text"
          class="w-full p-2 border rounded-lg"
          placeholder="Enter Scenario ID"
          required
        >
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700">Vimeo Video Link</label>
        <input
          v-model="videoLink"
          type="url"
          class="w-full p-2 border rounded-lg"
          placeholder="https://vimeo.com/example"
          required
        >
      </div>

      <button
        type="submit"
        class="bg-blue-500 text-white px-4 py-2 rounded"
      >
        Add Workbook
      </button>
    </form>

    <p
      v-if="message"
      class="mt-4 text-green-600"
    >
      {{ message }}
    </p>
  </div>
</template>

<script>
import { addWorkbook } from "@/services/workbookService";

export default {
  data() {
    return {
      scenarioId: "",
      videoLink: "",
      message: ""
    };
  },
  methods: {
    async createWorkbook() {
      try {
        await addWorkbook(this.scenarioId, this.videoLink);
        this.message = "Workbook successfully added!";
        this.scenarioId = "";
        this.videoLink = "";
      } catch (error) {
        console.error("Failed to add workbook:", error);
        this.message = "Error adding workbook.";
      }
    }
  }
};
</script>