<script lang="ts">
	import { api, type LeagueSummary } from '$lib/api';
	import { goto } from '$app/navigation';
	import { Users } from '@lucide/svelte';

	let leagues = $state<LeagueSummary[]>([]);
	let loaded = $state(false);
	let newName = $state('');
	let joinCode = $state('');
	let error = $state('');
	let busy = $state(false);

	async function load() {
		try {
			leagues = (await api.myLeagues()).leagues;
		} catch {
			/* ignore */
		} finally {
			loaded = true;
		}
	}
	$effect(() => {
		load();
	});

	async function create(e: Event) {
		e.preventDefault();
		error = '';
		busy = true;
		try {
			const r = await api.createLeague(newName);
			newName = '';
			goto(`/leagues/${r.id}`);
		} catch {
			error = 'No se pudo crear la liga.';
		} finally {
			busy = false;
		}
	}

	async function join(e: Event) {
		e.preventDefault();
		error = '';
		busy = true;
		try {
			const r = await api.joinLeague(joinCode);
			joinCode = '';
			goto(`/leagues/${r.id}`);
		} catch {
			error = 'Código de invitación inválido.';
		} finally {
			busy = false;
		}
	}
</script>

<p class="kicker">Juega con tus amigos</p>
<h1>Ligas</h1>
<p class="muted">Competencias privadas — tus pronósticos vs. los de tus amigos.</p>

<section class="card">
	<h3>Tus ligas</h3>
	{#if !loaded}
		<p class="muted">Cargando…</p>
	{:else if leagues.length === 0}
		<p class="muted">Ninguna aún — crea una o únete con un código.</p>
	{:else}
		{#each leagues as l (l.id)}
			<a class="lrow" href={`/leagues/${l.id}`}>
				<span>{l.name}</span>
				{#if l.role === 'owner'}<span class="pill">propietario</span>{/if}
				<span class="spacer"></span>
				<span class="cnt"><Users size={15} /> {l.members}</span>
			</a>
		{/each}
	{/if}
</section>

<section class="card">
	<h3>Crear una liga</h3>
	<form onsubmit={create}>
		<div class="field">
			<input class="input" placeholder="Nombre de la liga" bind:value={newName} required />
		</div>
		<button class="btn" disabled={busy || !newName.trim()}>Crear</button>
	</form>
</section>

<section class="card">
	<h3>Unirse a una liga</h3>
	<form onsubmit={join}>
		<div class="field">
			<input
				class="input code"
				placeholder="CÓDIGO DE INVITACIÓN"
				bind:value={joinCode}
				required
			/>
		</div>
		<button class="btn secondary" disabled={busy || !joinCode.trim()}>Unirse</button>
	</form>
</section>

{#if error}<p class="error">{error}</p>{/if}

<style>
	h1 {
		margin: 1rem 0 0.2rem;
	}
	.muted {
		margin: 0 0 1rem;
	}
	.lrow {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.7rem 0;
		border-top: 1px solid var(--border);
		color: var(--text);
	}
	.lrow:first-of-type {
		border-top: none;
	}
	.cnt {
		display: inline-flex;
		align-items: center;
		gap: 0.3rem;
		color: var(--muted);
		font-size: 0.9rem;
	}
	.code {
		text-transform: uppercase;
		letter-spacing: 0.2em;
		font-weight: 700;
	}
</style>
