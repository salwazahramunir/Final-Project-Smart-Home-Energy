import React, { useState } from "react";
import axios from "axios";

function App() {
  const [file, setFile] = useState(null);
  const [question, setQuestion] = useState("");
  const [query, setQuery] = useState("");
  const [response, setResponse] = useState("");

  const [isListening, setIsListening] = useState(false);
  const [language, setLanguage] = useState('en-US'); // Default ke bahasa Inggris

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const handleUpload = async () => {
    const formData = new FormData();
    formData.append("file", file);
    formData.append('question', question);

    try {
      const res = await axios.post('http://localhost:8080/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      setResponse(res.data.answer); // Assuming the response has an 'answer' field
    } catch (error) {
      console.error('Error uploading file:', error);
    }
  };

  const handleChat = async () => {
    try {
      const res = await axios.post("http://localhost:8080/chat", { query });
      setResponse(res.data.answer);
    } catch (error) {
      console.error("Error querying chat:", error);
    }
  };

  const SpeechRecognition =
    window.SpeechRecognition || window.webkitSpeechRecognition;
  const recognition = new SpeechRecognition();

  recognition.continuous = true;
  recognition.interimResults = true;
  recognition.lang = language; // Bahasa yang dipilih

  const handleStart = () => {
    setIsListening(true);
    recognition.start();

    recognition.onresult = (event) => {
      const transcript = Array.from(event.results)
        .map((result) => result[0].transcript)
        .join('');
      setQuery(transcript);
    };

    recognition.onerror = (event) => {
      console.error('Error occurred:', event.error);
    };
  };

  const handleStop = () => {
    setIsListening(false);
    recognition.stop();
  };

  const handleLanguageChange = (event) => {
    setLanguage(event.target.value);
  };

  return (
    <div style={{ maxWidth: "600px", margin: "0 auto", padding: "20px", textAlign: "center", fontFamily: "Arial, sans-serif" }}>
      <h1 style={{ color: "#333", marginBottom: "20px" }}>Data Analysis Chatbot</h1>
      <div style={{ marginBottom: "20px" }}>
        <input type="file" onChange={handleFileChange} style={{ padding: "10px", marginRight: "10px", border: "1px solid #ccc", borderRadius: "4px" }} />
        <input type="text" value={question} onChange={(e) => setQuestion(e.target.value)} style={{ padding: "10px", marginRight: "10px", border: "1px solid #ccc", borderRadius: "4px", height: "20px" }} />
        <button onClick={handleUpload} style={{ padding: "10px 20px", backgroundColor: "#007bff", color: "white", border: "none", borderRadius: "4px", cursor: "pointer" }}>
          Upload and Analyze
        </button>
      </div>
      <div style={{ marginBottom: "20px" }}>
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Ask a question..."
          style={{ padding: "10px", marginRight: "10px", border: "1px solid #ccc", borderRadius: "4px", width: "calc(100% - 140px)" }}
        />
        <button onClick={handleChat} style={{ padding: "10px 20px", backgroundColor: "#007bff", color: "white", border: "none", borderRadius: "4px", cursor: "pointer" }}>
          Chat
        </button>
      </div>
      <div style={{ marginBottom: '10px' }}>
        <label htmlFor="language">Select Language: </label>
        <select
          id="language"
          value={language}
          onChange={handleLanguageChange}
          style={{ padding: '5px' }}
        >
          <option value="en-US">English (US)</option>
          <option value="id-ID">Bahasa Indonesia</option>
          <option value="fr-FR">French</option>
          <option value="es-ES">Spanish</option>
          <option value="zh-CN">Mandarin (Simplified)</option>
          {/* Tambahkan bahasa lain sesuai kebutuhan */}
        </select>
      </div>
      <div>
        {isListening ? (
          <button onClick={handleStop}>Stop Listening</button>
        ) : (
          <button onClick={handleStart}>Start Listening</button>
        )}
      </div>
      <div style={{ marginTop: "20px", padding: "10px", border: "1px solid #ccc", borderRadius: "4px", backgroundColor: "#f9f9f9" }}>
        <h2>Response</h2>
        <p>{response}</p>
      </div>
    </div>
  );
}

export default App;