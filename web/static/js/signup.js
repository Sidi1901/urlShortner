async function signup() {
  const name            = document.getElementById('name').value.trim();
  const email           = document.getElementById('email').value.trim();
  const password        = document.getElementById('password').value;
  const confirmPassword = document.getElementById('confirmPassword').value;
  const userType        = document.getElementById('userType').value;

  // ── Client-side validation ──────────────────────────────────────────────
  if (!name || !email || !password || !confirmPassword) {
    showError('All fields are required.');
    return;
  }

  if (password !== confirmPassword) {
    showError('Passwords do not match.');
    return;
  }

  if (password.length < 8) {
    showError('Password must be at least 8 characters.');
    return;
  }

  // ── Call the backend API ────────────────────────────────────────────────
  // Endpoint: POST /users/signup
  // Required body fields: email, name, password, user_type ("free" | "premium")
  const btn = document.getElementById('submit-btn');
  btn.disabled = true;
  btn.textContent = 'Creating account…';
  hideError();

  try {
    const response = await fetch('/users/signup', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name:      name,
        email:     email,
        password:  password,
        user_type: userType,
      }),
    });

    const data = await response.json();

    if (response.status === 201) {
      // Success — show modal
      document.getElementById('modal').style.display = 'flex';
      return;
    }

    // Map known backend error codes to friendly messages
    if (response.status === 409) {
      showError('An account with this email already exists.');
    } else if (response.status === 400) {
      showError(data.error || 'Invalid input. Please check your details.');
    } else {
      showError(data.error || 'Something went wrong. Please try again.');
    }

  } catch (err) {
    showError('Network error. Please check your connection and try again.');
  } finally {
    btn.disabled = false;
    btn.textContent = 'Sign Up';
  }
}

function goToLogin() {
  window.location.href = '/login';
}

function showError(message) {
  const banner = document.getElementById('error-banner');
  banner.textContent = message;
  banner.style.display = 'block';
}

function hideError() {
  document.getElementById('error-banner').style.display = 'none';
}
