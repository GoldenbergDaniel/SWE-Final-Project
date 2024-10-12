<script>
  let showLogin = false; // Track whether to show the login form
  let showSignUp = false; // Track whether to show the sign-up form

  const handleSignInClick = () => {
    showLogin = true; 
    showSignUp = false; // Hide sign-up form when login is shown
  };

  const handleSignUpClick = () => {
    showSignUp = true;
    showLogin = false; // Hide login form when sign-up is shown
  };

  const handleBackClick = () => {
    showLogin = false;
    showSignUp = false; // Reset to the initial state
  };
</script>

<main>
  {#if !showLogin && !showSignUp}
    <!-- Initial State: Show Sign In and Sign Up Buttons -->
    <button on:click={handleSignInClick}>Sign In</button>
    <button on:click={handleSignUpClick}>Sign Up</button>
  
  {:else if showLogin}
    <!-- Login Form -->
    <div>
      <h2>Login</h2>
      <label>
        Username:
        <input type="text" placeholder="Enter your username" />
      </label>
      <label>
        Password:
        <input type="password" placeholder="Enter your password" />
      </label>
      <button>Submit</button>
      <button on:click={handleBackClick}>Back</button> <!-- Add a back button -->
    </div>
  
  {:else if showSignUp}
    <!-- Sign Up Form -->
    <div>
      <h2>Sign Up</h2>
      <label>
        Username:
        <input type="text" placeholder="Enter your username" />
      </label>
      <label>
        Password:
        <input type="password" placeholder="Enter your password" />
      </label>
      <button>Submit</button>
      <button on:click={handleBackClick}>Back</button> <!-- Add a back button -->
    </div>
  {/if}
</main>

<style>
  button {
    padding: 10px 20px;
    font-size: 16px;
    margin: 10px 0;
  }

  label {
    display: block;
    margin: 10px 0;
  }

  input {
    padding: 8px;
    font-size: 14px;
    margin-top: 5px;
  }

  div {
    margin-top: 20px;
  }
</style>
