import { db } from "@/firebase";
import { collection, addDoc } from "firebase/firestore";

// Function to add a new workbook
export const addWorkbook = async (scenarioId, videoLink) => {
  try {
    const workbookData = {
      scenarioId,
      videoLink,
      overview: {
        background: "Description of the event...",
        logSources: ["Firewall", "Endpoint Detection", "Authentication Logs"],
        commonality: "Rare"
      },
      investigation: {
        validation: {
          initialAlert: "Suspicious login detected",
          verificationSteps: [
            "Check authentication logs",
            "Confirm geolocation of login",
            "Was the response automated?"
          ]
        },
        scoping: {
          artifacts: ["User: jdoe", "IP: 192.168.1.1", "File Hash: abc123"],
          spreadAnalysis: [
            "Check system logs for further activity",
            "Identify impacted accounts"
          ]
        }
      },
      mitigation: {
        steps: ["Revoke compromised credentials", "Isolate infected endpoints"],
        COOP: "Ensure VPN access remains available",
        recovery: "Rebuild compromised machines"
      },
      review: {
        infrastructureChanges: "Enforce MFA",
        processImprovements: "Improve alert triaging",
        conclusion: "Lesson learned: Early detection is critical"
      }
    };

    const docRef = await addDoc(collection(db, "workbooks"), workbookData);
    console.log("Workbook created with ID:", docRef.id);
    return docRef.id;
  } catch (error) {
    console.error("Error adding workbook:", error);
    throw error;
  }
};