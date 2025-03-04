import { initializeApp } from 'firebase/app';
import { getAuth } from 'firebase/auth';
import { getFirestore } from 'firebase/firestore';
import { getStorage } from "firebase/storage";

const firebaseConfig = {
  apiKey: "AIzaSyDD7lHwzdzhIPs6mXZC5TjRJDreYELtwk8",
  authDomain: "replaydata-385e9.firebaseapp.com",
  projectId: "replaydata-385e9",
  storageBucket: "replaydata-385e9.appspot.com",
  messagingSenderId: "493455164147",
  appId: "1:493455164147:web:d6a3b2e041213a4a92ed86"
};
//  replaydata-385e9.firebasestorage.app
// Initialize Firebase
const app = initializeApp(firebaseConfig);
const auth = getAuth(app);
const db = getFirestore(app);
const storage = getStorage(app);

export { auth, db, storage };
