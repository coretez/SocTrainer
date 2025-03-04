<template>
    <div class="container mx-auto p-6">
      <h1 class="text-3xl font-bold text-gray-800 mb-8">Create New Scenario</h1>
  
      <div v-if="loading" class="text-lg text-gray-600 italic">Loading Profile...</div>
  
      <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div class="bg-white shadow-md rounded-lg overflow-hidden">
          <div class="px-6 py-4 bg-gray-50 border-b border-gray-200">
            <h2 class="text-xl font-semibold text-gray-700">Profile Information</h2>
          </div>
          <div class="p-6 space-y-3">
            <p class="text-gray-700">
              <span class="font-medium">Search Site Token:</span>
              {{ profile.ScenarioSiteToken }}
            </p>
            <p class="text-gray-700">
              <span class="font-medium">Base Search URL:</span>
              {{ profile.ScenarioSiteUrl }}
            </p>
          </div>
        </div>
  
        <div class="bg-white shadow-md rounded-lg overflow-hidden">
          <div class="px-6 py-4 bg-gray-50 border-b border-gray-200">
            <h2 class="text-xl font-semibold text-gray-700">Scenario Details</h2>
          </div>
          <div class="p-6 space-y-4">
            <label class="block text-gray-700 text-sm font-bold mb-2" for="scenarioName">
              Scenario Name
            </label>
            <input
              v-model="newScenario.name"
              type="text"
              id="scenarioName"
              placeholder="Enter scenario name"
              class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            />
  
            <label class="block text-gray-700 text-sm font-bold mb-2" for="scenarioDescription">
              Scenario Description
            </label>
            <input
              v-model="newScenario.description"
              type="text"
              id="scenarioDescription"
              placeholder="Enter scenario description"
              class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            />
            <label class="block text-gray-700 text-sm font-bold mb-2" for="gridAccount">
              Grid Account
            </label>
            <input
              v-model="searchData.gridAccount"
              type="text"
              id="gridAccount"
              placeholder="Enter grid account name"
              class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            />
  
            <label class="block text-gray-700 text-sm font-bold mb-2" for="searchString">
              Search String
            </label>
            <input
              v-model="searchData.searchString"
              type="text"
              id="searchString"
              placeholder="Enter search string"
              class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            />
  
            <label class="block text-gray-700 text-sm font-bold mb-2" for="startTime">
              Start Time (milliseconds)
            </label>
            <input
              v-model="searchData.startTime"
              type="number"
              id="startTime"
              placeholder="Enter start time in milliseconds"
              class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            />
  
            <label class="block text-gray-700 text-sm font-bold mb-2" for="endTime">
              End Time (milliseconds)
            </label>
            <input
              v-model="searchData.endTime"
              type="number"
              id="endTime"
              placeholder="Enter end time in milliseconds"
              class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            />
  
            <button
              @click="getData"
              class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
            >
              Get Data
            </button>
  
            <button
              @click="createScenario"
              class="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
            >
              Create Scenario
            </button>
            <!-- Data Display Area -->
            <div v-if="apiData" class="mt-4">
              <h3 class="text-lg font-semibold text-gray-700">API Data</h3>
              <pre class="bg-gray-100 rounded p-3 overflow-auto">
                {{ JSON.stringify(apiData, null, 2) }}
              </pre>
            </div>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import { ref, onMounted } from "vue";
  import { getFirestore, doc, getDoc, setDoc, collection } from "firebase/firestore";
  import { getAuth, onAuthStateChanged } from "firebase/auth";
  import { useRouter } from "vue-router";
  import axios from "axios"; // Import Axios
  
  export default {
    name: "AddScenario",
    setup() {
      const profile = ref({});
      const loading = ref(true);
      const apiData = ref(null);
      const searchData = ref({
        searchString: "",
        startTime: null,
        endTime: null,
        gridAccount: "", //ADD this!
        site:"",
        token:""
      });
      const newScenario = ref({ name: "", description: "" });
      const filename = ref({ name: "somename.json" });
      const user = ref(null);
      const router = useRouter();
  
      const db = getFirestore();
      const auth = getAuth();
  
      const fetchProfile = async () => {
        if (!user.value) {
          console.warn("🔴 User not authenticated. Skipping fetch.");
          loading.value = false;
          return;
        }
  
        try {
          const docRef = doc(db, "profiles", user.value.uid);
          const docSnap = await getDoc(docRef);
  
          if (docSnap.exists()) {
            profile.value = docSnap.data();
            console.log("✅ Profile data:", profile.value);
          } else {
            console.warn("⚠️ No profile found for this user!");
            profile.value = { message: "No profile created" };
          }
        } catch (error) {
          console.error("🔥 Error fetching profile:", error);
        } finally {
          loading.value = false;
        }
      };
  
      const getData = async () => {
             console.log(profile.value.ScenarioSiteUrl)  //If these values show up, it works
             console.log(profile.value.ScenarioSiteToken )
             console.log( searchData.value )
  
          //Add data from profiles to searchdata
               searchData.value.site=profile.value.ScenarioSiteUrl
               searchData.value.token=profile.value.ScenarioSiteToken
  
  
              console.log('Getting data with:', searchData.value);
              try {
                  // Make a POST request to your Go backend's /api/get-data endpoint
                  const response = await axios.post('http://localhost:8080/api/get-data', searchData.value);
  
                  // Handle the response from the backend
                  console.log('✅ Data from backend:', response.data);
              //   alert(`Created this to a client: ${response.data}`)  //If you get a message, delete code!
  
                console.log("Stringifying the data")
                 const stringData = JSON.stringify(response.data)
                 alert(`String  ${stringData}`)   //Set this up for now, delete once we're done with step 5
  
                  //Store from server
                   apiData.value = response.data;  //Set it to work
  
  
  
              } catch (error) {
                  console.error('🔥 Error fetching data from backend:', error);
                  alert(`🔥 Error fetching data from backend: ${error.message}`);
              }
  
          }
  
          const createScenario = async () => {
              try {
                  // Create a new document reference in the "scenarios" collection
                  const newScenarioRef = doc(collection(db, "scenarios"));
  
                  // Set the data for the new scenario
                  await setDoc(newScenarioRef, {
                      ScenarioName: newScenario.name,  //Use Scenario name
                      ScenarioData: searchData.value,    //Search Data
                    BackendData: apiData.value,  //add api data
                      description: newScenario.description,   // Use description
                      FileName: filename.name, //Hardcode filename
                  });
  
                  console.log("✅ Created a new scenario with:", newScenario.name);
                  alert(`Created a new scenario named ${newScenario.name} with ${newScenario.description}`);
                  router.push('/scenarios'); // Redirect to the scenarios page
              } catch (error) {
                  console.error("🔥 Error adding the scenario:", error);
                  alert(`Error creating the scenario: ${error.message}`);
              }
          };
  
          onMounted(() => {
              onAuthStateChanged(auth, (loggedInUser) => {
                  if (loggedInUser) {
                      console.log(`✅ User authenticated: ${loggedInUser.uid}`);
                      user.value = loggedInUser;
                      fetchProfile();
                  } else {
                      console.warn("❌ No user logged in.");
                      loading.value = false;
                  }
              });
          });
  
          return { profile, loading, searchData, newScenario,filename, getData, createScenario, apiData };
      },
  };
  </script>
  
  <style scoped>
  /* You may not need any scoped styles */
  </style>