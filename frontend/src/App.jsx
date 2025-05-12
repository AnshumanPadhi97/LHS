import { useState } from "react";
import { BuildStack, RunStackById } from "../wailsjs/go/main/App";

function App() {
  const [yaml, setYaml] = useState("");
  const [stackName, setStackName] = useState("");
  const [stackId, setStackId] = useState("");
  const [log, setLog] = useState("");

  const handleBuild = async () => {
    setLog("Building stack...");
    try {
      await BuildStack(yaml);
      const name = extractStackName(yaml);
      setStackName(name);
      setLog(`✅ Stack '${name}' built successfully.`);
    } catch (err) {
      setLog("❌ " + (err).message);
    }
  };

  const handleRun = async () => {
    if (!stackId) {
      setLog("❌ Please enter a stack ID.");
      return;
    }
    setLog("Running stack...");
    try {
      await RunStackById(parseInt(stackId));
      setLog(`✅ Stack with ID '${stackId}' is now running.`);
    } catch (err) {
      setLog("❌ " + (err).message);
    }
  };

  const extractStackName = (yamlContent) => {
    const match = yamlContent.match(/name:\s*(\w[\w-]*)/);
    return match ? match[1] : "unknown";
  };

  return (
    <div style={{ padding: 20 }}>
      <h2>🚀 Local Hosting Services</h2>
      <textarea
        rows={15}
        cols={80}
        placeholder="Paste your stack YAML here"
        value={yaml}
        onChange={(e) => setYaml(e.target.value)}
        style={{ fontFamily: "monospace", width: "100%", marginBottom: 10 }}
      />
      <button onClick={handleBuild}>Build Stack</button>
      <div style={{ marginTop: 10 }}>
        <input
          type="number"
          placeholder="Enter Stack ID to run"
          value={stackId}
          onChange={(e) => setStackId(e.target.value)}
          style={{ padding: 8, fontSize: 14, width: "100%" }}
        />
      </div>
      <button onClick={handleRun} style={{ marginTop: 10 }}>
        Run Stack
      </button>
      <pre style={{ background: "#111", color: "#0f0", padding: 10, marginTop: 20 }}>
        {log}
      </pre>
    </div>
  );
}

export default App;
