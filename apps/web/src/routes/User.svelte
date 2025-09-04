<script>
  import { updatePassword } from "../lib/controllers/users";

  /** @type {{id: string, name: string}}*/
  let { id, name } = $props();
  let loadingUpdatePassword = $state(false);
  let currentPassword = $state("");
  let newPassword = $state("");
  let repeatPassword = $state("");
  let error = $state({});

  const validatePasswordUpdateFields = (current, future, repeat) => {
    let dirtied = false;
    const error = {};
    if (current.length === 0) {
      error["current-password"] = "current password is required";
      dirtied = true;
    }

    if (future !== repeat) {
      error["repeat-password"] = "repeat password does not match new pasword";
      dirtied = true;
    }

    const newPasswordErrors = [];
    if (current === future) {
      newPasswordErrors.push("current password is equal to new password");
      dirtied = true;
    }
    if (future.length === 0) {
      newPasswordErrors.push("new password is required");
      dirtied = true;
    }

    if (future.length < 8) {
      newPasswordErrors.push("new password length should be longer than 8");
      dirtied = true;
    }

    if (newPasswordErrors.length > 0) {
      error["new-password"] = newPasswordErrors;
      dirtied = true;
    }

    if (dirtied) {
      return error;
    }

    return null;

    // TODO: add more password restrictions here (symbol, uppercase, lowercase, number)
  };

  /**
   * @param {SubmitEvent} event
   */
  const handleUpdatePasswordSubmit = async (event) => {
    event.preventDefault();

    error = {};
    const e = validatePasswordUpdateFields(
      currentPassword,
      newPassword,
      repeatPassword,
    );
    if (e) {
      error = e;
      return;
    }

    loadingUpdatePassword = true;
    try {
      await updatePassword({
        oldPassword: currentPassword,
        newPassword,
        repeatPassword,
      });

      currentPassword = "";
      newPassword = "";
      repeatPassword = "";
    } catch (e) {
      console.error(e);
    } finally {
      loadingUpdatePassword = false;
    }
  };
</script>

<div class="container">
  <h1 class="title is-1">{name}</h1>
  <div class="section">
    <h2 class="title is-2">Update password</h2>
    <form onsubmit={handleUpdatePasswordSubmit}>
      <div class="field">
        <label class="label" for="current-password">Current password</label>
        <div class="control">
          <input
            class={`input ${error["current-password"] ? "is-danger" : ""}`}
            type="password"
            id="current-password"
            bind:value={currentPassword}
            placeholder="Current password"
            disabled={loadingUpdatePassword}
          />
        </div>
        {#if error["current-password"]}
          <p class="help is-danger">{error["current-password"]}</p>
        {/if}
      </div>

      <div class="field">
        <label class="label" for="new-password">New password</label>
        <div class="control">
          <input
            class={`input ${error["new-password"] ? "is-danger" : ""}`}
            type="password"
            id="new-password"
            bind:value={newPassword}
            placeholder="New password"
            disabled={loadingUpdatePassword}
          />
        </div>
        {#if error["new-password"]}
          {#each error["new-password"] as e}
            <p class="help is-danger">{e}</p>
          {/each}
        {/if}
      </div>

      <div class="field">
        <label class="label" for="repeat-password">Repeat password</label>
        <div class="control">
          <input
            class={`input ${error["repeat-password"] ? "is-danger" : ""}`}
            type="password"
            id="repeat-password"
            bind:value={repeatPassword}
            placeholder="Repeat password"
            disabled={loadingUpdatePassword}
          />
        </div>
        {#if error["repeat-password"]}
          <p class="help is-danger">{error["repeat-password"]}</p>
        {/if}
      </div>

      <div class="field">
        <div class="control">
          <button class="button is-primary" disabled={loadingUpdatePassword}
            >Update password</button
          >
        </div>
      </div>
    </form>
  </div>
</div>
