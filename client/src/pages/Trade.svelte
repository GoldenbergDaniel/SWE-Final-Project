<script>
    import Footer from '../lib/Footer.svelte';
    import Header from '../lib/Header.svelte';
    import Background from '../lib/Background.svelte';
    import { checkAuth } from '../assets/auth.js';
    import { onMount } from 'svelte';
    import { navigate } from "svelte-routing";

    onMount(() => {
        if (!checkAuth()) {
            alert('Access Denied! Login Required');
            navigate('/');
        }
    })

    let userBalance = 0;

    onMount(async () => {
        try
        {
            const response = await fetch("http://localhost:5174/userdata", {
                method: "GET",
                credentials: "include",
            })

            if (!response.ok)
            {
                throw new Error("Failed to fetch user data")
            }

            const data = await response.json()
            userBalance = data.balance
        }
        catch (error)
        {
            console.error("Error fetching user data: ", error)
        }
    })

    // Variables for handling input fields
    let ticker = '';
    let shareNumber = 0;  // Ensure it's initialized to 0
    let rationale = '';

    // Handle share number input validation
    function handleShareNumberInput(event) {
        // Ensure the value is a positive number (no negative values allowed)
        if (shareNumber < 0) {
            shareNumber = 0;
        }
    }

    // Handle "Place Order" button click
    function placeOrder() {
        // Here, you'd typically add logic to process the order
        console.log("Order placed for:", ticker, shareNumber, rationale);
        // Reset form fields (optional)
        ticker = '';
        shareNumber = 0;
        rationale = '';
    }

</script>

<main>
    <Background />
    <Header />
    <div class="content">
        <h1>Trade Page</h1>
    </div>
    <!-- Right side: Balance and Trade form -->
    <div class="trade">
        <h2>Your Balance</h2>
        <table>
            <thead>
                <tr>
                    <th>Balance</th>
                </tr>
            </thead>
            <tbody>
                <tr>
                    <td>${userBalance}</td>
                </tr>
            </tbody>
        </table>

        <h2>Make a Trade</h2>
        <form>
            <!-- Ticker input -->
            <input
                type="text"
                placeholder="TICKER"
                bind:value={ticker}
                class="input-field"
            />
            <!-- Share Number input -->
            <input
                type="number"
                placeholder="Share Number"
                bind:value={shareNumber}
                min="0"
                on:input={handleShareNumberInput}
                class="input-field"
            />
            <!-- Rationale input -->
            <input
                type="text"
                placeholder="Rationale"
                bind:value={rationale}
                class="input-field"
            />
            <!-- Place Order button -->
            <button type="button" on:click={placeOrder} class="place-order-btn">Place Order</button>
        </form>
    </div>
    <Footer />
</main>

<style>
    .content{
        margin-top: 96px;
    }
    /* Right side: Balance and Trade form */
    .trade {
        flex: 0 0 35%; /* Trade section takes less space */
        min-width: 250px;
        padding: 20px;
        border: 1px solid #ddd;
        background-color: #f9f9f9;
        border-radius: 8px;
    }

    .trade h2 {
        text-align: center;
        margin-bottom: 20px;
    }

    /* Balance table */
    .trade table {
        width: 100%;
        border-collapse: collapse;
        margin-bottom: 20px;
    }

    .trade th, .trade td {
        padding: 10px;
        text-align: center;
        border: 1px solid #ddd;
    }

    .trade th {
        background-color: #f4f4f4;
        font-weight: bold;
    }

    .trade td {
        font-size: 1.2em;
    }

    /* Input fields */
    .input-field {
        width: 100%;
        padding: 10px;
        margin-bottom: 15px;
        border: 1px solid #ddd;
        border-radius: 5px;
        font-size: 16px;
        box-sizing: border-box;
    }

    .input-field::placeholder {
        color: #aaa;
    }

    .input-field:focus {
        border-color: #4CAF50;
        outline: none;
    }

    /* Place Order button */
    .place-order-btn {
        width: 100%;
        padding: 12px;
        background-color: #4CAF50; /* Green */
        color: white;
        border: none;
        border-radius: 5px;
        font-size: 16px;
        cursor: pointer;
        transition: background-color 0.3s;
    }

    .place-order-btn:hover {
        background-color: #45a049;
    }

    .place-order-btn:active {
        background-color: #3e8e41;
    }
</style>