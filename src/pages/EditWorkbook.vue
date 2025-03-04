<template>
  <div class="p-6 max-w-3xl mx-auto">
    <h1 class="text-2xl font-bold mb-4">
      Edit Workbook
    </h1>

    <div
      v-if="loading"
      class="text-gray-500"
    >
      Loading workbook...
    </div>
    <div
      v-else-if="!workbook"
      class="text-red-500"
    >
      Workbook not found.
    </div>

    <div
      v-else
      class="space-y-6"
    >
      <!-- Scenario Title -->
      <h2 class="text-xl font-bold text-indigo-600">
        {{ workbook.scenarioName || "Unknown Scenario" }}
      </h2>

      <!-- Editable Fields -->
      <div
        v-for="(field, key) in editors"
        :key="key"
        class="mt-4"
      >
        <label class="block text-sm font-medium mb-1">{{ workbookFields[key].label }}</label>

        <!-- Toolbar -->
        <div class="mb-2 flex flex-wrap gap-1">
          <button
            :class="buttonClass(editors[key].isActive('bold'))"
            @click="editors[key].chain().focus().toggleBold().run()"
          >
            B
          </button>
          <button
            :class="buttonClass(editors[key].isActive('italic'))"
            @click="editors[key].chain().focus().toggleItalic().run()"
          >
            I
          </button>
          <button
            :class="buttonClass(editors[key].isActive('code'))"
            @click="editors[key].chain().focus().toggleCode().run()"
          >
            { }
          </button>
          <button
            :class="buttonClass(editors[key].isActive('heading', { level: 1 }))"
            @click="editors[key].chain().focus().toggleHeading({ level: 1 }).run()"
          >
            H1
          </button>
          <button
            :class="buttonClass(editors[key].isActive('heading', { level: 2 }))"
            @click="editors[key].chain().focus().toggleHeading({ level: 2 }).run()"
          >
            H2
          </button>
          <button
            :class="buttonClass(editors[key].isActive('heading', { level: 3 }))"
            @click="editors[key].chain().focus().toggleHeading({ level: 3 }).run()"
          >
            H3
          </button>
          <button
            :class="buttonClass(editors[key].isActive('bulletList'))"
            @click="editors[key].chain().focus().toggleBulletList().run()"
          >
            • List
          </button>
          <button
            :class="buttonClass(editors[key].isActive('orderedList'))"
            @click="editors[key].chain().focus().toggleOrderedList().run()"
          >
            1. List
          </button>
          <button
            :class="buttonClass(editors[key].isActive('codeBlock'))"
            @click="editors[key].chain().focus().toggleCodeBlock().run()"
          >
            `Code`
          </button>
          <button
            :class="buttonClass(false)"
            @click="editors[key].chain().focus().setHorizontalRule().run()"
          >
            ───
          </button>
          <button
            :class="buttonClass(false)"
            @click="editors[key].chain().focus().setHardBreak().run()"
          >
            ↵
          </button>
        </div>

        <!-- Editor UI (No Blue Focus, Clickable Full Area) -->
        <div
          class="editor-container border p-2 rounded-lg min-h-[120px]"
          @click="focusEditor(key)"
        >
          <editor-content :editor="editors[key]" />
        </div>
      </div>

      <!-- Save Button -->
      <button
        class="mt-4 bg-blue-500 text-white px-4 py-2 rounded"
        @click="saveWorkbook"
      >
        Save Changes
      </button>

      <!-- Back Button -->
      <router-link
        to="/workbook"
        class="mt-4 block text-sm text-gray-500 hover:text-gray-700"
      >
        Cancel
      </router-link>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onBeforeUnmount } from "vue";
import { useRoute, useRouter } from "vue-router";
import { doc, getDoc, updateDoc } from "firebase/firestore";
import { db } from "@/firebase";
import { Editor, EditorContent } from "@tiptap/vue-3";
import StarterKit from "@tiptap/starter-kit";
import Placeholder from "@tiptap/extension-placeholder";
import BulletList from "@tiptap/extension-bullet-list";
import OrderedList from "@tiptap/extension-ordered-list";
import ListItem from "@tiptap/extension-list-item";
import Heading from "@tiptap/extension-heading";

export default {
  components: { EditorContent },

  setup() {
    const route = useRoute();
    const router = useRouter();
    const workbook = ref(null);
    const loading = ref(true);
    const editors = ref({});

    const workbookFields = {
      scenarioIntroduction: { label: "Scenario Introduction" },
      validationStage: { label: "Validation Stage" },
      scopingStage: { label: "Scoping Stage" },
      recoveryResponseStage: { label: "Recovery & Response Stage" },
      hotwash: { label: "Hotwash" },
    };

    // Styling for toolbar buttons
    const buttonClass = (isActive) => [
      "px-2 py-1 text-xs rounded-md border transition",
      "bg-gray-200 hover:bg-gray-300",
      isActive ? "bg-blue-500 text-white" : "text-gray-700",
    ];

    // Focus editor when clicking anywhere in its container
    const focusEditor = (key) => {
      if (editors.value[key]) {
        editors.value[key].commands.focus();
      }
    };

    // Load Workbook Data
    const fetchWorkbook = async () => {
      try {
        const docRef = doc(db, "workbooks", route.params.id);
        const docSnap = await getDoc(docRef);
        
        if (docSnap.exists()) {
          workbook.value = docSnap.data();

          // Initialize editors with proper styling
          Object.keys(workbookFields).forEach((key) => {
            editors.value[key] = new Editor({
            extensions: [
              StarterKit,
              Placeholder.configure({
                placeholder: workbookFields[key].placeholder,
              }),
              BulletList,
              OrderedList,
              ListItem,
              Heading.configure({ levels: [1, 2, 3] }), // Ensure H1-H3 work
            ],
            content: workbook.value.fields[key] || "",
            editorProps: {
              attributes: {
                class: "prose focus:outline-none",
              },
            },
          });
          });

        } else {
          console.error("Workbook not found.");
        }
      } catch (error) {
        console.error("Error fetching workbook:", error);
      } finally {
        loading.value = false;
      }
    };

    // Save Workbook
    const saveWorkbook = async () => {
      if (!workbook.value) return;

      const updatedFields = Object.keys(workbookFields).reduce((acc, key) => {
        acc[key] = editors.value[key].getHTML();
        return acc;
      }, {});

      try {
        const docRef = doc(db, "workbooks", route.params.id);
        await updateDoc(docRef, { fields: updatedFields });
        alert("Workbook updated successfully!");
        router.push(`/workbook/${route.params.id}`);
      } catch (error) {
        console.error("Error updating workbook:", error);
        alert("Failed to update workbook.");
      }
    };

    onMounted(fetchWorkbook);

    onBeforeUnmount(() => {
      Object.values(editors.value).forEach((editor) => editor.destroy());
    });

    return { workbook, loading, editors, workbookFields, saveWorkbook, focusEditor, buttonClass };
  },
};
</script>

<style scoped>
/* Ensure the editor area is fully clickable */
.editor-container {
  @apply border p-3 rounded-lg min-h-[120px] cursor-text focus:ring-0 outline-none;
}

/* Remove blue focus outline */
.editor-container:focus-within {
  @apply ring-0 outline-none;
}
</style>