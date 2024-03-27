package mbbservice

// TODO: реализовать тесты сервиса

// func TestCreateEvent(t *testing.T) {
// 	service := NewEventsService(memorystorage.NewStorage())
// 	ctx := context.Background()

// 	event := model.Event{
// 		StartDatetime: time.Date(2024, 1, 5, 20, 0, 0, 0, time.Local),
// 		EndDatetime:   time.Date(2024, 1, 5, 20, 15, 0, 0, time.Local),
// 	}

// 	id, err := service.CreateEvent(ctx, event)
// 	event.Id = id

// 	require.NoError(t, err)
// 	require.Equal(t, 1, id)

// 	// Проверка, что событие было добавлено в хранилище
// 	storedEvent, err := service.GetEvent(ctx, id)
// 	require.NoError(t, err)
// 	require.Equal(t, event, storedEvent)

// 	// Проверка ошибки пересечений
// 	_, err = service.CreateEvent(ctx, event)
// 	require.Error(t, err)
// 	require.ErrorIs(t, err, errors.ErrDateTimeBusy)

// 	// Проверка ошибки некорректного периода
// 	_, err = service.CreateEvent(ctx, model.Event{
// 		StartDatetime: time.Date(2024, 1, 5, 20, 15, 0, 0, time.Local),
// 		EndDatetime:   time.Date(2024, 1, 5, 20, 0, 0, 0, time.Local),
// 	})
// 	require.Error(t, err)
// 	require.ErrorIs(t, err, errors.ErrIncorrectDates)
// }
