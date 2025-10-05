## 🧭 **1. Penamaan Branch**

```
<tipe>/<deskripsi-singkat-dipisahkan-dengan-strip>
```

### 🔹 **Tipe branch umum:**

| Tipe        | Kegunaan                                             | Contoh                      |
| ----------- | ---------------------------------------------------- | --------------------------- |
| `feature/`  | Untuk fitur baru                                     | `feature/login-page`        |
| `bugfix/`   | Untuk memperbaiki bug non-produksi                   | `bugfix/fix-login-error`    |
| `hotfix/`   | Untuk memperbaiki bug di produksi                    | `hotfix/payment-timeout`    |
| `release/`  | Untuk persiapan rilis versi baru                     | `release/v1.2.0`            |
| `chore/`    | Untuk tugas non-fungsional (misal update dependensi) | `chore/update-dependencies` |
| `refactor/` | Untuk perbaikan struktur kode tanpa ubah fungsi      | `refactor/auth-service`     |
| `test/`     | Untuk penambahan atau perbaikan test                 | `test/add-user-tests`       |
| `docs/`     | Untuk dokumentasi                                    | `docs/api-endpoints`        |

### 🔹 **Tips tambahan:**

* Gunakan **huruf kecil semua** dan **pisahkan dengan tanda minus (-)**.
* Hindari nama terlalu panjang (>5 kata).
* Kalau pakai issue tracker (misal Jira/GitHub Issues), bisa tambahkan ID:

  ```
  feature/123-login-auth
  bugfix/456-api-timeout
  ```

---

## 🧾 **2. Penamaan Commit Message**

Gunakan format **konvensi commit** seperti [Conventional Commits](https://www.conventionalcommits.org):

```
<type>(<scope>): <deskripsi singkat>
```

### 🔹 **Contoh umum:**

| Type       | Kegunaan                                    | Contoh                                           |
| ---------- | ------------------------------------------- | ------------------------------------------------ |
| `feat`     | Menambahkan fitur baru                      | `feat(auth): add JWT-based login`                |
| `fix`      | Memperbaiki bug                             | `fix(api): correct null pointer in user service` |
| `docs`     | Perubahan dokumentasi                       | `docs(readme): add installation steps`           |
| `style`    | Perubahan gaya kode (formatting, indentasi) | `style: format code using prettier`              |
| `refactor` | Refactor tanpa ubah fungsi                  | `refactor(auth): simplify token validation`      |
| `test`     | Menambahkan atau memperbaiki test           | `test(user): add unit test for profile update`   |
| `chore`    | Tugas maintenance                           | `chore(deps): update express to v5.0.0`          |
| `perf`     | Perbaikan performa                          | `perf(query): optimize database lookup`          |
| `build`    | Perubahan pada sistem build                 | `build(ci): add github actions workflow`         |

### 🔹 **Contoh commit message bagus:**

```
feat(user): implement profile picture upload
fix(login): handle invalid credentials properly
refactor(api): move auth middleware to separate file
docs: update API usage section
```

### 🔹 **Tambahan opsional (body & footer):**

```
fix(api): handle null response in user endpoint

Previously, the API crashed when the response was null.
Now it returns an empty object instead.

Closes #123
```

---

## 🧩 **3. Contoh Praktik Nyata**

### Workflow:

```bash
# Buat branch baru untuk fitur login
git checkout -b feature/login-page

# Kerjakan fitur, lalu commit
git add .
git commit -m "feat(login): add login form and validation"

# Push ke remote
git push origin feature/login-page
```

---

## 💡 **Ringkasan Cepat**

| Elemen     | Format                       | Contoh                                   |
| ---------- | ---------------------------- | ---------------------------------------- |
| **Branch** | `<type>/<short-description>` | `feature/register-api`                   |
| **Commit** | `<type>(<scope>): <message>` | `fix(auth): correct token refresh logic` |
