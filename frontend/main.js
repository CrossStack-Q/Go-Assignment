let selectedFile = null;

// Handle file selection
document.getElementById('file-input').addEventListener('change', function(e) {
    selectedFile = e.target.files[0];
    if (selectedFile) {
        document.getElementById('selected-file').textContent = 
            `Selected file: ${selectedFile.name} (${(selectedFile.size/1024).toFixed(2)} KB)`;
    }
});

// Handle file upload
async function uploadFile() {
    if (!selectedFile) {
        showStatus('Please select a file first!', 'error');
        return;
    }

    const formData = new FormData();
    formData.append('csv', selectedFile);

    try {
        const response = await fetch('http://localhost:8080/v1/upload', {
            method: 'POST',
            body: formData
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const result = await response.json();
        showStatus(`Upload successful: ${result.message}`, 'success');
    } catch (error) {
        showStatus(`Upload failed: ${error.message}`, 'error');
    }
}

function showStatus(message, type) {
    const statusDiv = document.getElementById('status');
    statusDiv.textContent = message;
    statusDiv.style.color = type === 'error' ? '#dc3545' : '#28a745';
}
