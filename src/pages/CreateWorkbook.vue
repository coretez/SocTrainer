<template>
  <div class="p-6 max-w-3xl mx-auto">
    <h1 class="text-2xl font-bold mb-4">Create Workbook</h1>

    <!-- Select Scenario -->
    <div>
      <label class="block text-sm font-medium mb-2">Select Scenario</label>
      <select v-model="selectedScenario" class="w-full p-2 border rounded-lg">
        <option v-for="scenario in availableScenarios" :key="scenario.id" :value="scenario">
          {{ scenario.name }}
        </option>
      </select>
      <p v-if="selectedScenario" class="text-gray-500 mt-2">{{ selectedScenario.description }}</p>
    </div>

    <!-- Editor Sections -->
    <div v-for="(field, key) in workbookFields" :key="key" class="mt-4">
      <label class="block text-sm font-medium mb-1">{{ field.label }}</label>

      <div v-if="editors[key]" class="border p-2 rounded-lg min-h-[120px] relative">
        <!-- Editor Toolbar -->
        <div class="mb-2 flex flex-wrap gap-1">
          <button @click="editors[key].chain().focus().toggleBold().run()" :class="buttonClass(editors[key].isActive('bold'))">B</button>
          <button @click="editors[key].chain().focus().toggleItalic().run()" :class="buttonClass(editors[key].isActive('italic'))">I</button>
          <button @click="editors[key].chain().focus().toggleCode().run()" :class="buttonClass(editors[key].isActive('code'))">{ }</button>
          <button @click="editors[key].chain().focus().toggleHeading({ level: 1 }).run()" :class="buttonClass(editors[key].isActive('heading', { level: 1 }))">H1</button>
          <button @click="editors[key].chain().focus().toggleHeading({ level: 2 }).run()" :class="buttonClass(editors[key].isActive('heading', { level: 2 }))">H2</button>
          <button @click="editors[key].chain().focus().toggleHeading({ level: 3 }).run()" :class="buttonClass(editors[key].isActive('heading', { level: 3 }))">H3</button>
          <button @click="editors[key].chain().focus().toggleBulletList().run()" :class="buttonClass(editors[key].isActive('bulletList'))">• List</button>
          <button @click="editors[key].chain().focus().toggleOrderedList().run()" :class="buttonClass(editors[key].isActive('orderedList'))">1. List</button>
          <button @click="editors[key].chain().focus().toggleCodeBlock().run()" :class="buttonClass(editors[key].isActive('codeBlock'))">`Code`</button>
          <button @click="editors[key].chain().focus().setHorizontalRule().run()" :class="buttonClass(false)">───</button>
          <button @click="editors[key].chain().focus().setHardBreak().run()" :class="buttonClass(false)">↵</button>
        </div>

        <!-- Placeholder -->
        <div v-if="editors[key].isEmpty" class="absolute top-12 left-4 text-gray-400 pointer-events-none">
          {{ field.placeholder }}
        </div>

        <!-- Editor Content -->
        <editor-content :editor="editors[key]" class="border p-2 min-h-[100px]" />
      </div>
    </div>

    <!-- Save Button -->
    <button @click="saveWorkbook" class="mt-4 bg-blue-500 text-white px-4 py-2 rounded">
      Save Workbook
    </button>
  </div>
</template>

<script>
import { ref, onMounted, onBeforeUnmount, nextTick } from "vue";
import { collection, getDocs, addDoc } from "firebase/firestore";
import { db, auth } from "@/firebase";
import { Editor, EditorContent } from "@tiptap/vue-3";
import StarterKit from "@tiptap/starter-kit";
import Placeholder from "@tiptap/extension-placeholder";

export default {
  components: { EditorContent },

  setup() {
    const availableScenarios = ref([]);
    const selectedScenario = ref(null);
    const editors = ref({});

    const buttonClass = (isActive) => [
      "px-2 py-1 text-xs rounded-md border",
      "bg-gray-200 hover:bg-gray-300 transition",
      isActive ? "bg-blue-500 text-white" : "text-gray-700",
    ];
    
    const workbookFields = {
      scenarioIntroduction: { 
        label: "Scenario Introduction",
        placeholder: "Describe the overall scenario and context for the exercise."
      },
      validationStage: { 
        label: "Validation Stage",
        placeholder: "Explain how you would validate whether the alert is legitimate."
      },
      scopingStage: { 
        label: "Scoping Stage",
        placeholder: "Describe how you would determine the impact of the event."
      },
      recoveryResponseStage: { 
        label: "Recovery & Response Stage",
        placeholder: "Detail the remediation steps and response plan."
      },
      hotwash: { 
        label: "Hotwash",
        placeholder: "Summarize lessons learned and after-action insights."
      },
    };

    // Initialize Tiptap Editors with Placeholder Support
    onMounted(async () => {
      await nextTick();

      Object.keys(workbookFields).forEach((key) => {
        editors.value[key] = new Editor({
          extensions: [
            StarterKit,
            Placeholder.configure({
              placeholder: workbookFields[key].placeholder, 
            }),
          ],
          content: "", 
        });
      });

      // Fetch available scenarios
      const querySnapshot = await getDocs(collection(db, "scenarios"));
      availableScenarios.value = querySnapshot.docs.map(doc => ({
        id: doc.id, 
        name: doc.data().name, 
        description: doc.data().description
      }));
    });

    // Cleanup editors before component unmounts
    onBeforeUnmount(() => {
      Object.values(editors.value).forEach(editor => editor.destroy());
    });

    // Save Workbook
    const saveWorkbook = async () => {
      if (!selectedScenario.value) {
        alert("Please select a scenario.");
        return;
      }

      const user = auth.currentUser;
      if (!user) {
        alert("You must be logged in to save.");
        return;
      }

      const workbookData = {
        scenarioId: selectedScenario.value.id,
        scenarioName: selectedScenario.value.name, // ✅ Store scenario name
        createdBy: user.uid,
        fields: Object.keys(workbookFields).reduce((acc, key) => {
          acc[key] = editors.value[key].getHTML();
          return acc;
        }, {}),
        createdAt: new Date(),
      };

      await addDoc(collection(db, "workbooks"), workbookData);
      alert("Workbook saved successfully!");
    };

    return { availableScenarios, selectedScenario, editors, workbookFields, buttonClass, saveWorkbook };
  },
};
</script>