
function login() {
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    if (!email || !password) {
    alert('Email and Password are required');
    return;
    }

    document.getElementById('modal').style.display = 'flex';
}

function closeModal() {
    document.getElementById('modal').style.display = 'none';
    window.location.href = '/';
}