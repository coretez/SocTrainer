import { createRouter, createWebHistory } from "vue-router";
import { getAuth, onAuthStateChanged } from "firebase/auth";

// Import pages
import Home from "@/pages/Home.vue";
import ManageScenarios from '@/pages/ManageScenarios.vue';
import AddScenario from '@/pages/AddScenario.vue';
import Profile from "@/pages/Profile.vue";
import Workbook from "@/pages/Workbook.vue";
import WorkbookDetails from "@/pages/WorkbookDetails.vue";
import EditWorkbook from "@/pages/EditWorkbook.vue";
import Login from "@/pages/Login.vue";
import CreateWorkbook from "@/pages/CreateWorkbook.vue";

const routes = [
  { path: "/", name: "Home", component: Home },
  { path: "/scenarios", name: "ManageScenarios", component: ManageScenarios, meta: { requiresAuth: true } },
  { path: "/add-scenario", name: "AddScenario", component: AddScenario, meta: { requiresAuth: true } },
  { path: "/profile", name: "Profile", component: Profile, meta: { requiresAuth: true } },
  { path: "/workbook", name: "Workbook", component: Workbook, meta: { requiresAuth: true } },
  { path: "/workbook/:id", name: "WorkbookDetails", component: WorkbookDetails, meta: { requiresAuth: true } },
  { path: "/create-workbook", name: "CreateWorkbook", component: CreateWorkbook, meta: { requiresAuth: true } },
  { path: "/edit-workbook/:id", name: "EditWorkbook", component: EditWorkbook, meta: { requiresAuth: true } },
  { path: "/login", name: "Login", component: Login },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// 🔹 Navigation Guard for Authentication
router.beforeEach((to, from, next) => {
  const auth = getAuth();
  
  // Ensure we check Firebase authentication state
  onAuthStateChanged(auth, (user) => {
    const isAuthenticated = user || localStorage.getItem("authenticated");

    if (to.matched.some((record) => record.meta.requiresAuth) && !isAuthenticated) {
      next("/login");
    } else {
      next();
    }
  });
});

export default router;