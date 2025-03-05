<template>
    <div class="mt-6 bg-white shadow-md rounded-lg p-6">
      <h3 class="text-lg font-semibold text-gray-700 mb-4">Find & Replace</h3>
  
      <div v-for="(item, index) in replaceData" :key="index" class="flex space-x-4 mb-3">
        <input v-model="item.find" type="text" placeholder="Find" class="input-style w-1/2" />
        <input v-model="item.replace" type="text" placeholder="Replace" class="input-style w-1/2" />
        <button @click="removeReplacePair(index)" class="btn btn-red">X</button>
      </div>
  
      <button @click="addReplacePair" class="btn btn-gray">+ Add Find/Replace</button>
      <button @click="applyReplacements" class="btn btn-yellow mt-2">Replace</button>
    </div>
  </template>
  
  <script>
  import { ref } from "vue";
  
  export default {
    name: "FindReplace",
    props: {
      apiData: Object, // Accept API data as a prop
    },
    emits: ["update-api-data"], // Emit changes back to parent
    setup(props, { emit }) {
      const replaceData = ref([]);
  
      const addReplacePair = () => replaceData.value.push({ find: "", replace: "" });
  
      const removeReplacePair = (index) => replaceData.value.splice(index, 1);
  
      const applyReplacements = () => {
        if (!props.apiData) return;
  
        let modifiedData = JSON.stringify(props.apiData);
        replaceData.value.forEach(({ find, replace }) => {
          if (find) {
            const regex = new RegExp(find, "g");
            modifiedData = modifiedData.replace(regex, replace);
          }
        });
  
        emit("update-api-data", JSON.parse(modifiedData)); // Send updated data back
      };
  
      return { replaceData, addReplacePair, removeReplacePair, applyReplacements };
    },
  };
  </script>
  
  <style scoped>
  .input-style { @apply shadow border rounded w-full py-2 px-3 text-gray-700; }
  .btn { @apply py-2 px-4 mx-2 rounded font-bold focus:outline-none; }
  .btn-blue { @apply bg-blue-500 hover:bg-blue-700 text-white; }
  .btn-green { @apply bg-green-500 hover:bg-green-700 text-white; }
  .btn-red { @apply bg-red-500 hover:bg-red-700 text-white; }
  .btn-gray { @apply bg-gray-500 hover:bg-gray-700 text-white; }
  .btn-yellow { @apply bg-yellow-500 hover:bg-yellow-700 text-white; }
  </style>