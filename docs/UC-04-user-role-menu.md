# UC-04 — User, Role & Menu Management

> **Status**: Active  
> **Owner**: Wilie / Fadhil  
> **Last Updated**: 2026-06-05  
> **Route Group**: `/v1/management/{user,role,menu,userlog}` (Auth-gated)

---

## Overview

Mengelola akun user (admin/operator), role (permission groups), menu (navigasi sidebar), dan log aktivitas user. Auth menggunakan JWT HS256.

---

## Endpoints

### User Management

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| POST | `/v1/management/user/` | `CreateUser` | Buat user baru |
| GET | `/v1/management/user/` | `GetUserTable` | List semua user (paginated) |
| GET | `/v1/management/user/usercounts` | `GetUserCounts` | Jumlah user total/aktif |
| PUT | `/v1/management/user/:id` | `UpdateUser` | Update profil user |
| PUT | `/v1/management/user/assignservice/:id` | `AssignService` | Assign service & adnet ke user |
| PUT | `/v1/management/user/updatestatus/:id` | `UpdateUserStatus` | Aktifkan/nonaktifkan user |
| DELETE | `/v1/management/user/:id` | `DeleteUser` | Hapus user |
| GET | `/v1/management/user/approvalrequest` | `GetUserApplovalRequestTable` | List user pending approval |
| PUT | `/v1/management/user/approveuser/:id` | `ApproveUser` | Approve user registration |

### Role Management

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| POST | `/v1/management/role/` | `CreateRole` | Buat role baru |
| GET | `/v1/management/role/` | `GetRoleTable` | List semua role |
| PUT | `/v1/management/role/:id` | `UpdateRole` | Update role |
| DELETE | `/v1/management/role/:id` | `DeleteRole` | Hapus role |

### Menu Management

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| POST | `/v1/management/menu/` | `CreateMenu` | Buat menu baru |
| GET | `/v1/management/menu/` | `GetAllMenus` | List semua menu |
| GET | `/v1/management/menu/:id` | `GetMenuByID` | Detail menu by ID |
| PUT | `/v1/management/menu/:id` | `UpdateMenu` | Update menu |
| DELETE | `/v1/management/menu/:id` | `DeleteMenu` | Hapus menu |

### User Log

| Method | Path | Handler | Deskripsi |
|--------|------|---------|-----------|
| GET | `/v1/management/userlog/ip` | `GetIpAddress` | Ambil IP address requester |
| POST | `/v1/management/userlog/` | `CreateUserLog` | Simpan log aktivitas |
| GET | `/v1/management/userlog/` | `DisplayUserLogList` | List semua log |
| GET | `/v1/management/userlog/:id` | `DisplayUserLogHistory` | History log per user |

---

## Auth Flow (JWT)

```
Request ke protected group
  │
  ├─ AUTH_ENFORCE_DEFAULT=true → AuthMiddleware dieksekusi
  │   ├─ Extract Bearer token dari header Authorization
  │   ├─ Verify JWT (HS256, claims: sub/jti/type/exp/nbf/aud/iss)
  │   ├─ Check blacklist: Redis key "jwt:blacklist:{jti}" (DB=REDISDBINDEX)
  │   └─ Inject claims ke c.Locals (companies, adnets, etc.)
  └─ Handler dieksekusi dengan context yang sudah authenticated
```

---

## Scope Filtering

Beberapa endpoint filtering secara otomatis berdasarkan scope dari JWT:
- `allowedCompanies` → dari `c.Locals("companies").([]string)`
- `allowedAdnets` → dari `c.Locals("adnets").([]string)`

Digunakan di report endpoints, bukan di management endpoints (management butuh full access).

---

## Handler Files
- [`incoming_user_management.go`](../src/interfaces/http/handler/incoming_user_management.go)
- [`incoming_role_management.go`](../src/interfaces/http/handler/incoming_role_management.go)
- [`incoming_menu_management.go`](../src/interfaces/http/handler/incoming_menu_management.go)
- [`incoming_user_log_handler.go`](../src/interfaces/http/handler/incoming_user_log_handler.go)
- [`routes/management.go`](../src/interfaces/http/routes/management.go)

---

## Menambah/Mengubah Fitur

1. Edit handler di file yang sesuai
2. Daftarkan route baru di `routes/management.go` → fungsi yang sesuai (`registerUser`, `registerRole`, dll)
3. Jika ada field baru → update entity + migrate
4. Untuk JWT claim baru → update `AuthMiddleware` + update tabel scope di sini
5. `go build ./...` → pastikan clean
