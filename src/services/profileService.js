import { db } from "@/firebase";
import { doc, getDoc, setDoc, updateDoc } from "firebase/firestore";
import { getAuth } from "firebase/auth";

// Fetch user profile
export const getUserProfile = async () => {
  const auth = getAuth();
  const user = auth.currentUser;
  if (!user) return null;

  const userRef = doc(db, "users", user.uid);
  const docSnap = await getDoc(userRef);

  if (docSnap.exists()) {
    return docSnap.data();
  } else {
    // If profile doesn't exist, create one
    const newProfile = {
      name: user.displayName || "",
      email: user.email || "",
      company: "",
      hecUrl: "",
      hecToken: "",
      profilePic: user.photoURL || ""
    };
    await setDoc(userRef, newProfile);
    return newProfile;
  }
};

// Update user profile
export const updateUserProfile = async (updatedProfile) => {
  const auth = getAuth();
  const user = auth.currentUser;
  if (!user) return;

  const userRef = doc(db, "users", user.uid);
  await updateDoc(userRef, updatedProfile);
};