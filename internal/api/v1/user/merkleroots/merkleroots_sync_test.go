package merkleroots_test

// func TestSyncMerkleRoots(t *testing.T) {
// 	t.Run("Should properly sync database when empty", func(t *testing.T) {
// 		// setup
// 		httpClient := resty.New()
// 		server := merklerootstest.MockMerkleRootsAPIResponseNormal()
// 		defer server.Close()

// 		xpriv, err := bip32.NewKeyFromString(clienttest.UserXPriv)
// 		require.NoError(t, err)

// 		// given
// 		repo := merklerootstest.CreateRepository([]models.MerkleRoot{})
// 		client, err := merkleroots.NewClient(server.URL, httpClient, xpriv)
// 		require.NoError(t, err)

// 		// when
// 		err = client.SyncMerkleRoots(context.Background(), repo)

// 		// then
// 		require.NoError(t, err)
// 		require.Len(t, repo.MerkleRoots, len(merklerootstest.MockedSPVWalletData))
// 		require.Equal(t, merklerootstest.LastMockedMerkleRoot(), repo.MerkleRoots[len(repo.MerkleRoots)-1])
// 	})

// 	t.Run("Should properly sync database when partially filled", func(t *testing.T) {
// 		// setup
// 		httpClient := resty.New()
// 		server := merklerootstest.MockMerkleRootsAPIResponseNormal()
// 		defer server.Close()

// 		xpriv, err := bip32.NewKeyFromString(clienttest.UserXPriv)
// 		require.NoError(t, err)

// 		// given
// 		repo := merklerootstest.CreateRepository([]models.MerkleRoot{
// 			{
// 				MerkleRoot:  "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b",
// 				BlockHeight: 0,
// 			},
// 			{
// 				MerkleRoot:  "0e3e2357e806b6cdb1f70b54c3a3a17b6714ee1f0e68bebb44a74b1efd512098",
// 				BlockHeight: 1,
// 			},
// 			{
// 				MerkleRoot:  "9b0fc92260312ce44e74ef369f5c66bbb85848f2eddd5a7a1cde251e54ccfdd5",
// 				BlockHeight: 2,
// 			},
// 		})
// 		client, err := merkleroots.NewClient(server.URL, httpClient, xpriv)
// 		require.NoError(t, err)

// 		// when
// 		err = client.SyncMerkleRoots(context.Background(), repo)

// 		// then
// 		require.NoError(t, err)
// 		require.Len(t, repo.MerkleRoots, len(merklerootstest.MockedSPVWalletData))
// 		require.Equal(t, merklerootstest.LastMockedMerkleRoot(), repo.MerkleRoots[len(repo.MerkleRoots)-1])
// 	})

// 	t.Run("Should fail sync merkleroots due to the timeout", func(t *testing.T) {
// 		// setup
// 		httpClient := resty.New()
// 		server := merklerootstest.MockMerkleRootsAPIResponseDelayed()
// 		defer server.Close()

// 		xpriv, err := bip32.NewKeyFromString(clienttest.UserXPriv)
// 		require.NoError(t, err)

// 		// given
// 		repo := merklerootstest.CreateRepository([]models.MerkleRoot{})
// 		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
// 		defer cancel()

// 		client, err := merkleroots.NewClient(server.URL, httpClient, xpriv)
// 		require.NoError(t, err)

// 		// when
// 		err = client.SyncMerkleRoots(ctx, repo)

// 		// then
// 		require.ErrorIs(t, err, goclienterr.ErrSyncMerkleRootsTimeout) // Match the correct error
// 	})

// 	t.Run("Should fail sync merkleroots due to last evaluated key being the same in the response", func(t *testing.T) {
// 		// setup
// 		httpClient := resty.New()
// 		server := merklerootstest.MockMerkleRootsAPIResponseStale()
// 		defer server.Close()

// 		xpriv, err := bip32.NewKeyFromString(clienttest.UserXPriv)
// 		require.NoError(t, err)

// 		// given
// 		repo := merklerootstest.CreateRepository([]models.MerkleRoot{})
// 		client, err := merkleroots.NewClient(server.URL, httpClient, xpriv)
// 		require.NoError(t, err)

// 		// when
// 		err = client.SyncMerkleRoots(context.Background(), repo)

// 		// then
// 		require.ErrorIs(t, err, errors.ErrStaleLastEvaluatedKey)
// 	})

// 	t.Run("Should fail when merkleRootsAPI is nil", func(t *testing.T) {
// 		// given
// 		repo := merklerootstest.CreateRepository([]models.MerkleRoot{})
// 		client := &merkleroots.Client{} // Uninitialized API

// 		// when
// 		err := client.SyncMerkleRoots(context.Background(), repo)

// 		// then
// 		require.EqualError(t, err, "merkleRootsAPI is not initialized")
// 	})
// }
