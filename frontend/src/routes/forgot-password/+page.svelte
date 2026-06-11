<script lang="ts">
	import { auth } from '$lib/auth.svelte';

	let email = $state('');
	let busy = $state(false);
	let sent = $state(false);
	let error = $state('');

	async function submit(e: Event) {
		e.preventDefault();
		error = '';
		busy = true;
		try {
			await auth.requestPasswordReset(email.trim());
			sent = true;
		} catch (err: unknown) {
			error =
				(err as { message?: string })?.message ??
				'No se pudo enviar el correo de recuperación.';
		} finally {
			busy = false;
		}
	}
</script>

<div class="auth">
	<h1>Restablecer contraseña</h1>
	<p class="muted">
		Ingresa el correo con el que te registraste y te enviaremos un enlace.
	</p>

	{#if sent}
		<div class="card">
			<p class="ok">Si ese correo está registrado, recibirás un enlace en breve.</p>
			<p class="muted switch"><a href="/login">Volver a iniciar sesión</a></p>
		</div>
	{:else}
		<form class="card" onsubmit={submit}>
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
			{#if error}<p class="error">{error}</p>{/if}
			<button class="btn" disabled={busy || !email.trim()}>
				{busy ? 'Enviando…' : 'Enviar enlace'}
			</button>
			<p class="muted switch"><a href="/login">Volver a iniciar sesión</a></p>
		</form>
	{/if}
</div>

<style>
	.auth {
		max-width: 380px;
		margin: 12dvh auto 0;
	}
	h1 {
		margin: 0;
		font-size: 1.8rem;
	}
	.muted {
		margin: 0.25rem 0 1.5rem;
	}
	.ok {
		color: var(--success);
		font-size: 0.95rem;
		margin: 0;
	}
	.switch {
		text-align: center;
		margin: 1rem 0 0;
	}
</style>
