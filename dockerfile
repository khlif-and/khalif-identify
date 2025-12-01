# --- STAGE 1: Builder ---
# Kita gunakan Image Go 1.24 (Versi Stabil di Docker Hub)
FROM golang:1.24-alpine AS builder

# Install git (Wajib untuk download module)
RUN apk add --no-cache git

# Set folder kerja di dalam container
WORKDIR /app

# Copy file dependency
COPY go.mod go.sum ./

# --- MAGIC FIX ---
# Perintah ini MEMAKSA file go.mod DI DALAM DOCKER turun ke versi 1.23
# Ini SOLUSI agar build tidak error "requires go >= 1.24",
# tanpa perlu mengubah file asli di laptop Anda.
RUN go mod edit -go=1.23

# Download dependency
RUN go mod download

# Copy source code sisanya
COPY . .

# Build aplikasi menjadi binary bernama 'server'
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

# --- STAGE 2: Runner ---
FROM alpine:latest

# Install sertifikat SSL (Wajib untuk Azure/HTTPS)
RUN apk --no-cache add ca-certificates tzdata

# Set Timezone (Opsional, agar log jamnya pas WIB)
ENV TZ=Asia/Jakarta

WORKDIR /root/

# Copy hasil build dari Stage 1
COPY --from=builder /app/server .

# Buat folder uploads
RUN mkdir -p uploads

# Expose port
EXPOSE 8081

# Jalankan aplikasi
CMD ["./server"]