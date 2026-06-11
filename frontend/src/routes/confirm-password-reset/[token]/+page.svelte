<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';

	let token = $derived($page.params.token ?? '');
	let password = $state('');
	let confirm = $state('');
	let busy = $state(false);
	let error = $state('');
	let done = $state(false);

	async function submit(e: Event) {
		e.preventDefault();
		error = '';
		if (password.length < 8) {
			error = 'La contraseña debe tener al menos 8 caracteres.';
			return;
		}
		if (password !== confirm) {
			error = 'Las contraseñas no coinciden.';
			return;
		}
		busy = true;
		try {
			await auth.confirmPasswordReset(token, password, confirm);
			done = true;
			// PocketBase invalidates the session after a reset — make sure
			// nothing stale lingers, then send the user to sign in fresh.
			auth.logout();
			setTimeout(() => goto('/login'), 1200);
		} catch (err: unknown) {
			error =
				(err as { message?: string })?.message ??
				'El enlace es inválido o ha expirado.';
		} finally {
			busy = false;
		}
	}
</script>

<div class="auth">
	<h1>Elige una nueva contraseña</h1>
	<p class="muted">Ingresa y confirma tu nueva contraseña.</p>

	{#if done}
		<div class="card">
			<p class="ok">Contraseña actualizada — redirigiendo…</p>
		</div>
	{:else}
		<form class="card" onsubmit={submit}>
			<div class="field">
				<label for="pw">Nueva contraseña</label>
				<input
					id="pw"
					class="input"
					type="password"
					bind:value={password}
					autocomplete="new-password"
					minlength="8"
					required
				/>
			</div>
			<div class="field">
				<label for="pw2">Confirmar contraseña</label>
				<input
					id="pw2"
					class="input"
					type="password"
					bind:value={confirm}
					autocomplete="new-password"
					minlength="8"
					required
				/>
			</div>
			{#if error}<p class="error">{error}</p>{/if}
			<button class="btn" disabled={busy || !token}>
				{busy ? 'Actualizando…' : 'Actualizar contraseña'}
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
