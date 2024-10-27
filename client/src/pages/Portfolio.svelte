<script>
    import Footer from '../lib/Footer.svelte';
    import Header from '../lib/Header.svelte';
    import Background from '../lib/Background.svelte';
</script>

<main>
    <Background />
    <Header />
    <div class="content">
        <h1>Portfolio Page</h1>
    </div>
    <Footer />
</main>

<style>
    .content{
        margin-top: 96px;
    }
</style>