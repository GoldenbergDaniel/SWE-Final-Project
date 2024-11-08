<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { navigate } from "svelte-routing";
  
  interface SignUpData {
    first_name: string;
    last_name: string;
    email: string;
    username: string;
    password: string;
  }
  
  interface LoginData {
    username: string;
    password: string;
  }
  
  export let showSignIn = false;
  export let showSignUp = false;
  let properLogin = true;
  let errorMessage = '';
  
  // Create form data bindings
  let formData: Partial<SignUpData> = {
    first_name: '',
    last_name: '',
    email: '',
    username: '',
    password: ''
  };
  
  const dispatch = createEventDispatcher();
  
  const API_BASE_URL = 'http://localhost:5174';

  const resetBodyHeight = () => {
    document.body.style.height = `auto`;
  };
  
  const handleSignInClick = () => {
    dispatch("updateSignIn", true);
    dispatch("updateSignUp", false);
    resetBodyHeight();
    resetForm();
  };
  
  const handleSignUpClick = () => {
    dispatch("updateSignUp", true);
    dispatch("updateSignIn", false);
    resetBodyHeight();
    resetForm();
  };
  
  const handleBackClick = () => {
    dispatch("updateSignIn", false);
    dispatch("updateSignUp", false);
    resetBodyHeight();
    resetForm();
  };
  
  const resetForm = () => {
    formData = {
      first_name: '',
      last_name: '',
      email: '',
      username: '',
      password: ''
    };
    properLogin = true;
    errorMessage = '';
  };
  
  async function makeAuthRequest(endpoint: string, data: LoginData | SignUpData) {
    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      credentials: "include",
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    });

    const responseText = await response.text();
    if (!response.ok) {
      throw new Error(responseText || 'Authentication failed');
    }

    return responseText;
  }
  
  const handleSubmitClick = async () => {
    try {
      const endpoint = showSignUp ? '/signup' : '/login';
      const userData = showSignUp 
        ? formData as SignUpData
        : { 
            username: formData.username, 
            password: formData.password 
          } as LoginData;

      const message = await makeAuthRequest(endpoint, userData);
      
      if (showSignUp) {
        showSignUp = false;
        showSignIn = true;
        alert('Sign up successful! Please log in.');
      } else {
        navigate("/dashboard");
      }
      
      resetForm();
    } catch (error) {
      console.error('Authentication failed:', error);
      properLogin = false;
      errorMessage = error instanceof Error ? error.message : 'Authentication failed';
    }
  }

  async function onEnter()
  {
    handleSubmitClick()
  }
</script>

<main>
  {#if !showSignIn && !showSignUp}
    <button class="sign-in" on:click={handleSignInClick}>Sign In</button>
    <button class="sign-up" on:click={handleSignUpClick}>Sign Up</button>

  {:else if showSignIn}
    <div>
      <h2>Sign In</h2>
      <label>
        Username
        <input type="text" name="username" bind:value={formData.username} placeholder="Enter your username"/>
      </label>
      <label>
        Password
        <input type="password" name="password" bind:value={formData.password} placeholder="Enter your password" />
      </label>
      {#if !properLogin}
        <p class="error-message">{errorMessage}</p>
      {/if}
      <div class="button-container">
        <button class="back" on:click={handleBackClick}>Back</button>
        <button class="submit" on:keydown|preventDefault={onEnter} on:click={handleSubmitClick}>Submit</button>
      </div>
    </div>

  {:else if showSignUp}
    <div>
      <h2>Sign Up</h2>
      <label>First Name: <input type="text" name="first_name" bind:value={formData.first_name} placeholder="Enter your name" /></label>
      <label>Last Name: <input type="text" name="last_name" bind:value={formData.last_name} placeholder="Enter your last name" /></label>
      <label>Email: <input type="email" name="email" bind:value={formData.email} placeholder="Enter your email" /></label>
      <label>Username: <input type="text" name="username" bind:value={formData.username} placeholder="Enter your username" /></label>
      <label>Password: <input type="password" name="password" bind:value={formData.password} placeholder="Enter your password" /></label>
      {#if !properLogin}
        <p class="error-message">{errorMessage}</p>
      {/if}
      <div class="button-container">
        <button class="back" on:click={handleBackClick}>Back</button>
        <button class="submit" on:click={handleSubmitClick}>Submit</button>
      </div>
    </div>
  {/if}
</main>

<style>
  div {
    margin-top: 1em;
    max-width: 400px;
    width: 100%;
    text-align: left;
  }

  label {
    display: block;
    margin: 0.5em 0;
  }

  input {
    padding: 0.6em;
    font-size: 0.9em;
    margin-top: 0.4em;
    width: 100%;
    box-sizing: border-box;
    color: white;
    border-color: antiquewhite;
    background-color: black;
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
    background-color: black;
    color: white;
  }

  .button-container {
    display: flex;
    justify-content: space-between;
  }

  .sign-in:hover, .submit:hover, .sign-up:hover {
    color: white;
    border: none;
    border-radius: 5px;
    background-color: black;
  }

  .sign-in:hover {
    background-color: #4CAF50;
  }

  .sign-up:hover {
    background-color: #c3112c;
  }

  .submit:hover {
    background-color: #4CAF50;
  }

  .back:hover {
    background-color: #c3112c;
    color: white;
    border: none;
    border-radius: 5px;
  }

  .error-message {
    color: red;
    font-size: 0.9em;
    margin-top: 0.5em;
  }
</style>
