import { getAuth, signInWithEmailAndPassword, signOut } from "firebase/auth";
import { db } from "@/firebase"; // Ensure firebase.js is set up

const auth = getAuth();

// Function to log in
export const login = async (email, password) => {
  try {
    const userCredential = await signInWithEmailAndPassword(auth, email, password);
    return userCredential.user;
  } catch (error) {
    throw new Error(error.message);
  }
};

// Function to log out
export const logout = async () => {
  await signOut(auth);
  localStorage.removeItem("authenticated");
};
