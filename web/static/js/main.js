/**
 * Lordis Backend Engine - Core Client Pipeline
 * Handled Features: Async Forms, State Validations, Custom Toasts, DOM UI Lifecycle
 */

document.addEventListener("DOMContentLoaded", () => {
  // ==========================================
  // 🗃️ 1. GLOBAL UI CORE LAYOUT SELECTORS
  // ==========================================
  const loginForm = document.querySelector("form[action='/login']");
  const registerForm = document.querySelector("form[action='/register']");
  const forgotPasswordForm = document.querySelector("form[action='/forgot-password']");
  const resetPasswordForm = document.querySelector("form[action='/reset-password']");
  const profileUpdateForm = document.getElementById("profile-update-form");
  const logoutBtn = document.getElementById("logout-btn");
  const sidebarToggle = document.getElementById("sidebar-toggle");
  const sidebarContainer = document.querySelector(".sidebar");

  // ==========================================
  // 🔐 2. AUTHENTICATION & FORM PIPELINES
  // ==========================================

  // 🔑 A. User Authentication Handler (POST /login)
  if (loginForm) {
    loginForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      if (!validateFormFields(loginForm)) return;

      const bodyData = new URLSearchParams(new FormData(loginForm));
      toggleLoadingState(loginForm, true);

      try {
        const response = await fetch("/login", {
          method: "POST",
          headers: { "Content-Type": "application/x-www-form-urlencoded" },
          body: bodyData,
        });

        if (response.ok) {
          showToast("Authentication confirmed. Entering engine...", "success");
          window.location.href = "/dashboard";
        } else {
          const errorMsg = await response.text();
          showToast(`Authentication Error: ${errorMsg || "Invalid credentials."}`, "error");
        }
      } catch (err) {
        console.error("Critical Engine Link Down:", err);
        showToast("Network disconnect: The Go server failed to respond.", "error");
      } finally {
        toggleLoadingState(loginForm, false);
      }
    });
  }

  // 📝 B. User Profile Registration Pipeline (POST /register)
  if (registerForm) {
    registerForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      if (!validateFormFields(registerForm)) return;

      const pass = registerForm.querySelector("input[name='password']")?.value;
      const confirmPass = registerForm.querySelector("input[name='confirm_password']")?.value;

      if (confirmPass !== undefined && pass !== confirmPass) {
        showToast("Input Discrepancy: Password fields do not match.", "error");
        return;
      }

      const bodyData = new URLSearchParams(new FormData(registerForm));
      toggleLoadingState(registerForm, true);

      try {
        const response = await fetch("/register", {
          method: "POST",
          headers: { "Content-Type": "application/x-www-form-urlencoded" },
          body: bodyData,
        });

        if (response.ok) {
          showToast("Profile written to database successfully! Redirecting...", "success");
          setTimeout(() => { window.location.href = "/login"; }, 1500);
        } else {
          const errorMsg = await response.text();
          showToast(`Registration Aborted: ${errorMsg || "Invalid registration data."}`, "error");
        }
      } catch (err) {
        console.error("Registration Engine Link Down:", err);
        showToast("Network fault: Registration transaction failed.", "error");
      } finally {
        toggleLoadingState(registerForm, false);
      }
    });
  }

  // 📬 C. Password Reset Request Flow (POST /forgot-password)
  if (forgotPasswordForm) {
    forgotPasswordForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      if (!validateFormFields(forgotPasswordForm)) return;

      const bodyData = new URLSearchParams(new FormData(forgotPasswordForm));
      toggleLoadingState(forgotPasswordForm, true);

      try {
        const response = await fetch("/forgot-password", {
          method: "POST",
          headers: { "Content-Type": "application/x-www-form-urlencoded" },
          body: bodyData,
        });

        if (response.ok) {
          showToast("Reset Token compiled! Check your email inbox.", "success");
          forgotPasswordForm.reset();
        } else {
          const errorMsg = await response.text();
          showToast(`Dispatch Error: ${errorMsg || "Could not issue reset request."}`, "error");
        }
      } catch (err) {
        showToast("Failed to route token dispatch sequence.", "error");
      } finally {
        toggleLoadingState(forgotPasswordForm, false);
      }
    });
  }

  // 🔄 D. Secure Password Overwrite (POST /reset-password)
  if (resetPasswordForm) {
    resetPasswordForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      if (!validateFormFields(resetPasswordForm)) return;

      const bodyData = new URLSearchParams(new FormData(resetPasswordForm));
      toggleLoadingState(resetPasswordForm, true);

      try {
        const response = await fetch("/reset-password", {
          method: "POST",
          headers: { "Content-Type": "application/x-www-form-urlencoded" },
          body: bodyData,
        });

        if (response.ok) {
          showToast("Password updated cleanly! Routing to authorization...", "success");
          setTimeout(() => { window.location.href = "/login"; }, 1500);
        } else {
          const errorMsg = await response.text();
          showToast(`Overwrite Blocked: ${errorMsg || "Reset token invalid or expired."}`, "error");
        }
      } catch (err) {
        showToast("Network error trying to rewrite password credentials.", "error");
      } finally {
        toggleLoadingState(resetPasswordForm, false);
      }
    });
  }

  // 🚪 E. Session De-authorization & Logout Execution
  if (logoutBtn) {
    logoutBtn.addEventListener("click", async (e) => {
      e.preventDefault();
      try {
        const response = await fetch("/logout", { method: "POST" });
        if (response.ok) {
          showToast("Session killed safely. Goodbye.", "success");
          window.location.href = "/login";
        } else {
          showToast("Failed to clear session attributes on the server.", "error");
        }
      } catch (err) {
        console.error("Logout runtime interrupt:", err);
      }
    });
  }

  // ==========================================
  // 📊 3. DYNAMIC DASHBOARD DATA AND WORKFLOWS
  // ==========================================

  // 🛠️ Profile Schema Core Field Modification
  if (profileUpdateForm) {
    profileUpdateForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      const bodyData = new URLSearchParams(new FormData(profileUpdateForm));

      try {
        const response = await fetch("/api/user/update", {
          method: "PUT",
          headers: { "Content-Type": "application/x-www-form-urlencoded" },
          body: bodyData,
        });

        if (response.ok) {
          showToast("Internal database records updated.", "success");
        } else {
          showToast("Engine rejected metadata updates.", "error");
        }
      } catch (err) {
        showToast("Failed to serialize database updates for profile parameters.", "error");
      }
    });
  }

  // ==========================================
  // 📱 4. VISUAL EFFECTS & SYSTEM HELPERS
  // ==========================================

  // Responsive Interface Sidebar Drawer Controller
  if (sidebarToggle && sidebarContainer) {
    sidebarToggle.addEventListener("click", (e) => {
      e.stopPropagation();
      sidebarContainer.classList.toggle("collapsed-state");
    });
  }

  // Baseline Form Real-time Validation Helper
  function validateFormFields(formElement) {
    const requiredInputs = formElement.querySelectorAll("input[required], select[required]");
    let validState = true;

    requiredInputs.forEach((input) => {
      if (!input.value.trim()) {
        input.style.borderColor = "#ef4444";
        validState = false;
      } else {
        input.style.borderColor = "#e2e8f0";
      }
    });

    if (!validState) {
      showToast("Please fill out all required fields.", "warning");
    }
    return validState;
  }

  // UI Feedback Submission States Control
  function toggleLoadingState(form, isLoading) {
    const primaryButton = form.querySelector("button[type='submit']");
    if (!primaryButton) return;

    if (isLoading) {
      primaryButton.disabled = true;
      primaryButton.dataset.originalText = primaryButton.innerHTML;
      primaryButton.innerHTML = `<span class="spinner-loader">Working...</span>`;
    } else {
      primaryButton.disabled = false;
      if (primaryButton.dataset.originalText) {
        primaryButton.innerHTML = primaryButton.dataset.originalText;
      }
    }
  }
});

// ==========================================
// 🔔 5. GLOBAL NOTIFICATION DISMISSAL & TOASTS
// ==========================================

// Global Dismiss Function (Supports direct onclick="dismissNotification('123')")
window.dismissNotification = async function (id) {
  const elem = document.getElementById(`notif-${id}`);
  if (elem) {
    elem.style.opacity = "0";
    elem.style.transform = "scale(0.95)";
    setTimeout(() => elem.remove(), 300);
  }

  try {
    await fetch(`/notifications/read/${id}`, { method: "POST" });
  } catch (err) {
    console.error("Failed to sync read status with Go backend:", err);
  }
};

// Pure JavaScript Custom Alert Engine (Toast System)
window.showToast = function (message, type = "success") {
  const toast = document.createElement("div");
  toast.className = `custom-toast toast-${type}`;
  toast.textContent = message;

  let bgColor = "#10b981"; // Success Green
  if (type === "error") bgColor = "#ef4444";   // Error Red
  if (type === "warning") bgColor = "#f59e0b"; // Warning Orange

  Object.assign(toast.style, {
    position: "fixed",
    bottom: "30px",
    right: "30px",
    padding: "14px 28px",
    backgroundColor: bgColor,
    color: "#ffffff",
    borderRadius: "8px",
    fontWeight: "500",
    boxShadow: "0 10px 15px -3px rgba(0,0,0,0.1), 0 4px 6px -2px rgba(0,0,0,0.05)",
    zIndex: "10000",
    fontFamily: "system-ui, -apple-system, sans-serif",
    transition: "all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275)",
    opacity: "0",
    transform: "translateY(20px)",
  });

  document.body.appendChild(toast);

  requestAnimationFrame(() => {
    toast.style.opacity = "1";
    toast.style.transform = "translateY(0)";
  });

  setTimeout(() => {
    toast.style.opacity = "0";
    toast.style.transform = "translateY(10px)";
    setTimeout(() => toast.remove(), 400);
  }, 3500);
};