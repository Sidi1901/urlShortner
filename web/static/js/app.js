
function setExpiry(hours) {
    const now = new Date();
    now.setHours(now.getHours() + hours);

    document.getElementById('expiryDate').value = now.toISOString().split('T')[0];
    document.getElementById('expiryTime').value = now.toTimeString().slice(0,5);
}

function generateUrl() {
    const longUrl = document.getElementById('longUrl').value;

    if (!longUrl) {
    alert('Long URL is required');
    return;
    }

    const shortUrl = 'http://localhost:8080/abc123';

    document.getElementById('shortLink').innerText = shortUrl;
    document.getElementById('shortLink').href = shortUrl;

    document.getElementById('modal').style.display = 'flex';
}

function closeModal() {
    document.getElementById('modal').style.display = 'none';
}
