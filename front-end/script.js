const currentDate = new Date();
document.getElementById('current-date').textContent = currentDate.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
});

const dueDate = new Date(currentDate);
dueDate.setDate(dueDate.getDate() + 30);
document.getElementById('due-date').textContent = dueDate.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
});

async function downloadPDF() {
    try {
        const content = document.getElementById('bill-content');
        const styles = document.querySelector('link[href="styless.css"]');

        const cssResponse = await fetch(styles.href);
        const css = await cssResponse.text();

        const container = document.createElement('div');
        const styleElement = document.createElement('style');
        styleElement.textContent = css;
        container.appendChild(styleElement);
        container.appendChild(content.cloneNode(true));

        const response = await fetch('/api/generate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ html: container.outerHTML })
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'invoice.pdf';
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        window.URL.revokeObjectURL(url);
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to generate PDF. Please try again.');
    }
}
