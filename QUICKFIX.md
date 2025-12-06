# QUICK FIX: Add Audience Mapper in Keycloak UI

## The Problem
The 401 error occurs because JWT tokens don't have the correct `aud` (audience) claim.

## The Solution (5 minutes)

### Step 1: Open Keycloak Admin Console
1. Go to: **http://localhost:8081**
2. Click **Administration Console**
3. Login:
   - Username: `admin`
   - Password: `Admin#1!`

### Step 2: Navigate to the Client
1. In the left sidebar, select **Realm: traveler-dev** (dropdown at top)
2. Click **Clients** in the left menu
3. Click on **traveler-app**

### Step 3: Add the Audience Mapper
1. Click the **Client scopes** tab
2. Click on **traveler-app-dedicated** (in the "Assigned client scope" table)
3. Click **Add mapper** → **By configuration**
4. Select **Audience**
5. Fill in the form:
   - **Name**: `audience-mapper`
   - **Included Client Audience**: `traveler-app`
   - **Add to ID token**: **OFF** (unchecked)
   - **Add to access token**: **ON** (checked)
6. Click **Save**

### Step 4: Test in Your IDE
1. Go back to GoLand
2. Open `api/offerings/specials.http`
3. Make sure **dev** environment is selected
4. Run request **#1** (Get token)
5. Run request **#2** (Call specials endpoint)

**Expected result**: ✅ HTTP 200 with JSON response

---

## Alternative: Automated Fix

If you prefer an automated approach:

```bash
cd /Users/ettienemare/GolandProjects/traveler
./scripts/recreate-keycloak.sh
```

⚠️ **Warning**: This will recreate the Keycloak container and reset any manual changes.

---

## Verification

After adding the mapper, you can verify the token has the correct claims:

1. Get a new token in your HTTP client (run request #1 again)
2. Copy the `access_token` value
3. Go to: https://jwt.io
4. Paste the token
5. Check the payload - you should see:
   ```json
   {
     "iss": "http://localhost:8081/realms/traveler-dev",
     "aud": "traveler-app",
     "azp": "traveler-app",
     ...
   }
   ```

The `aud` field should now be present!

---

## Still Not Working?

### Check these:

1. **Did you get a NEW token after adding the mapper?**
   - Old tokens won't have the new claim
   - Run request #1 again in your HTTP client

2. **Is the Traveler API running?**
   ```bash
   curl http://localhost:8080/api/ping
   ```
   Should return JSON with status: "ok"

3. **Check application logs:**
   ```bash
   tail -f logs/app.log
   ```
   Look for JWT validation error messages

4. **Verify config matches:**
   - File: `configs/config.yaml`
   - Should have:
     ```yaml
     auth:
       issuer: http://localhost:8081/realms/traveler-dev
       audience: traveler-app
     ```

---

## Screenshots (for reference)

### 1. Clients Page
Navigate to: Clients → traveler-app

### 2. Client Scopes Tab
Click the "Client scopes" tab, then click "traveler-app-dedicated"

### 3. Add Mapper
Click "Add mapper" → "By configuration" → "Audience"

### 4. Mapper Configuration
- Name: `audience-mapper`
- Included Client Audience: `traveler-app`
- Add to access token: ✓ ON

---

## Need More Help?

See the full troubleshooting guide:
- `docs/troubleshooting/README.md`
- `docs/troubleshooting/401-specials-endpoint.md`

