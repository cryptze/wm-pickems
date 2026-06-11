<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';

	// After registering, resume an invite if one was carried in the URL.
	let invite = $derived($page.url.searchParams.get('invite'));
	function dest() {
		return invite ? `/join/${invite}` : '/';
	}
	let loginHref = $derived(
		invite ? `/login?invite=${encodeURIComponent(invite)}` : '/login'
	);

	let name = $state('');
	let email = $state('');
	let password = $state('');
	let error = $state('');
	let busy = $state(false);

	async function submit(e: Event) {
		e.preventDefault();
		error = '';
		if (password.length < 8) {
			error = 'La contraseña debe tener al menos 8 caracteres.';
			return;
		}
		busy = true;
		try {
			await auth.register(name, email, password);
			goto(dest());
		} catch (err: unknown) {
			error =
				(err as { message?: string })?.message ??
				'No se pudo crear la cuenta.';
		} finally {
			busy = false;
		}
	}
</script>

<div class="auth">
	<h1>Crear cuenta</h1>
	<p class="muted">Únete a la quiniela del Mundial.</p>

	<form class="card" onsubmit={submit}>
		<div class="field">
			<label for="nm">Nombre</label>
			<input id="nm" class="input" bind:value={name} required />
		</div>
		<div class="field">
			<label for="em">Correo electrónico</label>
			<input
				id="em"
				class="input"
				type="email"
				bind:value={email}
				autocomplete="email"
				required
			/>
		</div>
		<div class="field">
			<label for="pw">Contraseña</label>
			<input
				id="pw"
				class="input"
				type="password"
				bind:value={password}
				autocomplete="new-password"
				required
			/>
		</div>
		{#if error}<p class="error">{error}</p>{/if}
		<button class="btn" disabled={busy}>{busy ? 'Creando…' : 'Crear cuenta'}</button>
		<p class="muted switch">
			¿Ya tienes cuenta? <a href={loginHref}>Inicia sesión</a>
		</p>
	</form>
</div>

<style>
	.auth {
		max-width: 380px;
		margin: 10dvh auto 0;
	}
	h1 {
		margin: 0;
		font-size: 1.8rem;
	}
	.muted {
		margin: 0.25rem 0 1.5rem;
	}
	.switch {
		text-align: center;
		margin: 1rem 0 0;
	}
</style>
