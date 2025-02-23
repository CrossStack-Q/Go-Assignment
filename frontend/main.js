let selectedFile = null;

document.getElementById('file-input').addEventListener('change', function (e) {
    selectedFile = e.target.files[0];
    if (selectedFile) {
        document.getElementById('selected-file').textContent =
            `Selected file: ${selectedFile.name} (${(selectedFile.size / 1024).toFixed(2)} KB)`;
    }
});

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
        
        
        const fileCode = result.output_file.split('/').pop().replace('_output.csv', '');
        showOutputFileCode(`File Code: ${fileCode}`, 'success');
    } catch (error) {
        showStatus(`Upload failed: ${error.message}`, 'error');
        showOutputFileCode(`File Code: Nil`, 'error');
    }
}

async function checkStatus() {
    const fileCode = document.getElementById('file-code').value.trim();
    if (!fileCode) {
        showStatusResult('Please enter a file code.', 'error');
        return;
    }

    try {
        const response = await fetch('http://localhost:8080/v1/checkStatus', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ code: fileCode }),
        });

        const result = await response.json();

        if (response.ok) {
            console.log(result)
            showStatusResult(`Status: ${result.status}`, 'success');
        } else {
            showStatusResult(`Status: Failure - ${result.message}`, 'error');
        }
    } catch (error) {
        showStatusResult(`Error checking status: ${error.message}`, 'error');
    }
}

function showStatus(message, type) {
    const statusDiv = document.getElementById('status');
    statusDiv.textContent = message;
    statusDiv.style.color = type === 'error' ? '#dc3545' : '#28a745';
}

function showOutputFileCode(message, type) {
    const outputDiv = document.getElementById('outputFile');
    outputDiv.textContent = message;
    outputDiv.style.color = type === 'error' ? '#dc3545' : '#28a745';
}

function showStatusResult(message, type) {
    const statusDiv = document.getElementById('statusResult');
    statusDiv.textContent = message;
    statusDiv.style.color = type === 'error' ? '#dc3545' : '#28a745';
}
