<script>
  import { createEventDispatcher } from "svelte";

  export let showSignIn = false;
  export let showSignUp = false;
  let properLogin = true;

  const dispatch = createEventDispatcher();

  const handleSignInClick = () => {
    dispatch("updateSignIn", true);  // Signal to show sign-in
    dispatch("updateSignUp", false); // Ensure sign-up is hidden
  };

  const handleSignUpClick = () => {
    dispatch("updateSignUp", true);  // Signal to show sign-up
    dispatch("updateSignIn", false); // Ensure sign-in is hidden
  };

  const handleBackClick = () => {
    dispatch("updateSignIn", false); // Reset both to false
    dispatch("updateSignUp", false);
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
    <button class="sign-in" on:click={handleSignInClick}>Sign In</button>
    <button class ="sign-up" on:click={handleSignUpClick}>Sign Up</button>
  
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
      <div class="button-container">
        <button class="submit" on:click={handleSubmitClick}>Submit</button>
        <button class="back" on:click={handleBackClick}>Back</button> <!-- Add a back button -->
      </div>
    </div>
  
  {:else if showSignUp}
    <!-- Sign Up Form -->
    <div>
      <h2>Sign Up</h2>
      <label>
        First Name:
        <input type="text" placeholder="Enter your name" />
      </label>
      <label>
        Last Name:
        <input type="text" placeholder="Enter your last name" />
      </label>
      <label>
        Email:
        <input type="email" placeholder="Enter your email" />
      </label>
      <label>
        Username:
        <input type="text" placeholder="Enter your username" />
      </label>
      <label>
        Password:
        <input type="password" placeholder="Enter your password" />
      </label>
      <div class="button-container">
        <button class="submit">
          {#if properLogin}
            <a href="../dashboard.html">Submit</a>
          {/if}
        </button>
        <button class="back" on:click={handleBackClick}>Back</button> <!-- Add a back button -->
      </div>
    </div>
  {/if}
</main>

<style>

  div {
    margin-top: 1em;
    max-width: 400px; /* Limit form width */
    width: 100%; /* Ensure it adapts to max-width */
    text-align: left; /* Align labels and inputs within the div */
  }

  a {
    color: white;
  }

  label {
    display: block;
    margin: 0.5em 0;
    text-align: left; /* Align the label text to the right */
  }

  input {
    padding: 0.6em;
    font-size: 0.9em;
    margin-top: 0.4em;
    width: 100%; /* Make all inputs the same width */
    box-sizing: border-box; /* Ensure padding doesn't affect width */
    color:white;
    border-color:antiquewhite;
  }

  input:focus {
    outline: none;
    box-shadow: 0 0 5px 2px #4CAF50;
    border-color: #4CAF50;
  }

  button {
    padding: 0.7em 1.2em;
    font-size: 1em;
    margin: 2em 0;
  }

  .button-container {
    display: flex;
    justify-content: space-between;
  }

  .sign-in, .sign-up, .submit, .back {
    z-index: 2;
  }

  .sign-in:hover {
    background-color: #4CAF50;
    color: white;
    border: none;
    border-radius: 5px;
  }

  .sign-in:focus{
    outline: none;
    box-shadow: 0 0 5px 2px #4CAF50;
    border-color: #4CAF50;
  }

  .sign-up:hover {
    background-color: #c3112c;
    color: white;
    border: none;
    border-radius: 5px;
  }

  .sign-up:focus{
    outline: none;
    box-shadow: 0 0 5px 2px #c3112c;
    border-color: #c3112c;
  }

  .submit:hover{
    background-color: #4CAF50;
    color: white;
    border: none;
    border-radius: 5px;
  }

  .back:hover{
    background-color: #c3112c;
    color: white;
    border: none;
    border-radius: 5px;
  }

  .submit:focus{
    outline: none;
    box-shadow: 0 0 5px 2px #4CAF50;
    border-color: #4CAF50;
  }

  .back:focus{
    outline: none;
    box-shadow: 0 0 5px 2px #c3112c;
    border-color: #c3112c;
  }

</style>