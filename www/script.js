document.addEventListener("DOMContentLoaded", () => {
    // Find all buttons with a data-script attribute
    const buttons = document.querySelectorAll("button[data-script]");
  
    buttons.forEach(button => {
      button.addEventListener("click", async () => {
        const scriptPath = button.getAttribute("data-script");
  
        try {
          // Send a POST request to the server to execute the script
          const response = await fetch("/execute-script", {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({ scriptPath }),
          });
  
          if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
          }
  
          const result = await response.json();
          alert(result.message); // Display the result from the server
        } catch (error) {
          console.error("Error executing script:", error);
          alert("An error occurred while executing the script.");
        }
      });
    });
  });