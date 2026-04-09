

function signup(e) {

    console.log("E1",e)

    e.preventDefault();

    const name = document.getElementById('name').value;
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirmPassword').value;

    if (!name || !email || !password || !confirmPassword) {
    alert('All fields are required');
    return;
    }

    if (password !== confirmPassword) {
    alert('Passwords do not match');
    return;
    }

    console.log("E3")

    document.getElementById('modal').style.display = 'flex';

    const signupForm = document.getElementById("signup-form")

    const formData = new FormData(signupForm)
    const data = Object.fromEntries(formData)

    console.log(data)
}

function closeModal() {
    document.getElementById('modal').style.display = 'none';
    // window.location.href = '/';
}
