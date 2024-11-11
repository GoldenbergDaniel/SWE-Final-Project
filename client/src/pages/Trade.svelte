<script>
    import Footer from '../lib/Footer.svelte';
    import Header from '../lib/Header.svelte';
    import Background from '../lib/Background.svelte';
    import { isAuthenticated } from '../assets/auth.js';
    import { onMount } from 'svelte';
    import { navigate } from "svelte-routing";

    onMount(() => {
        if (!$isAuthenticated) {
            alert('Access Denied! Login Required');
            navigate('/');
        }
    });

</script>

<main>
    <Background />
    <Header />
    <div class="content">
        <h1>Trade Page</h1>
    </div>
    <Footer />
</main>

<style>
    .content{
        margin-top: 96px;
    }
</style>