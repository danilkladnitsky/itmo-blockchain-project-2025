# Use specific Python 3.9 version
FROM python:3.9.18-slim

# Set working directory
WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Install core dependencies first
RUN pip install --no-cache-dir \
    numpy==1.24.3 \
    pandas==1.5.3 \
    scikit-learn==1.2.2

# Copy requirements first to leverage Docker cache
COPY requirements.txt .

# Install remaining dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy the rest of the application
COPY . .

# Create directory for models if it doesn't exist
RUN mkdir -p train

# Expose the port the app runs on
EXPOSE 8000

# Command to run the application
CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "8000"]
