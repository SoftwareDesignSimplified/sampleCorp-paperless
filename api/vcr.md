Database

1. Use the real database -> slow / unreliable
2. Mock the database -> fast / brittle -> coupling to implementation
3. In memory database / repository -> SQLite
4. Nullable architecture

External API

1. Use the real API -> slow / unreliable
2. Mock the API -> fast / brittle -> coupling to implementation
3. In memory API / repository
4. Nullable architecture
5. VCR

Application code -> Repository -> Database

Application code -> Repository -> Infrastructure Wrapper -> Database

Application code -> Repository -> Null Infrastructure Wrapper
