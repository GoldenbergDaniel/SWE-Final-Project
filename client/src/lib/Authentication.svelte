<script>
  import App from "../App.svelte";
  import { createEventDispatcher } from "svelte";

  export let showSignIn = false; // Track whether to show the login form
  export let showSignUp = false; // Track whether to show the sign-up form
  export let showMainPage = true;

  const dispatch = createEventDispatcher();

  const handleSignInClick = () => {
    showSignIn = true; 
    showSignUp = false; // Hide sign-up form when login is shown
  };

  const handleSignUpClick = () => {
    showSignUp = true;
    showSignIn = false; // Hide login form when sign-up is shown
  };

  const handleBackClick = () => {
    showSignIn = false;
    showSignUp = false; // Reset to the initial state
  };

  const handleSubmitClick = () => {
  const usernameInput = document.querySelector('input[type="text"]');
  const passwordInput = document.querySelector('input[type="password"]');

  const userData = {
    username: usernameInput,
    password: passwordInput
  };

  dispatch('userData', userData);
};

</script>

<main>
  {#if !showSignIn && !showSignUp}
    <!-- Initial State: Show Sign In and Sign Up Buttons -->
    <button on:click={handleSignInClick}>Sign In</button>
    <button on:click={handleSignUpClick}>Sign Up</button>
  
  {:else if showSignIn}
    <!-- Sign In Form -->
    <div>
      <h2>Sign In</h2>
      <label>
        Username
        <input type="text" placeholder="Enter your username"/>
      </label>
      <label>
        Password
        <input type="password" placeholder="Enter your password" />
      </label>
      <button on:click={handleSubmitClick}>Submit</button>
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
