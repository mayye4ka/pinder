package service

import (
	"pinder/models"
	"pinder/server"
)

func mapPreferences(pref server.Preferences) models.Preferences {
	return models.Preferences{
		MaxAge:           pref.MaxAge,
		MinAge:           pref.MinAge,
		Gender:           models.Gender(pref.Gender),
		LocationLat:      pref.LocationLat,
		LocationLon:      pref.LocationLon,
		LocationRadiusKm: pref.LocationRadiusKm,
	}
}

func mapProfile(prof server.Profile) models.Profile {
	return models.Profile{
		Name:         prof.Name,
		Gender:       models.Gender(prof.Gender),
		Age:          prof.Age,
		Bio:          prof.Bio,
		LocationLat:  prof.LocationLat,
		LocationLon:  prof.LocationLon,
		LocationName: prof.LocationName,
	}
}

func unmapPreferences(pref models.Preferences) server.Preferences {
	return server.Preferences{
		MaxAge:           pref.MaxAge,
		MinAge:           pref.MinAge,
		Gender:           server.Gender(pref.Gender),
		LocationLat:      pref.LocationLat,
		LocationLon:      pref.LocationLon,
		LocationRadiusKm: pref.LocationRadiusKm,
	}
}

func unmapProfile(prof models.Profile, photos []string) server.Profile {
	return server.Profile{
		Name:         prof.Name,
		Gender:       server.Gender(prof.Gender),
		Age:          prof.Age,
		Bio:          prof.Bio,
		Photos:       photos,
		LocationLat:  prof.LocationLat,
		LocationLon:  prof.LocationLon,
		LocationName: prof.LocationName,
	}
}

func unmapContentType(content models.MsgContentType) server.MsgContentType {
	if content == models.ContentPhoto {
		return server.ContentPhoto
	}
	if content == models.ContentText {
		return server.ContentText
	}
	return server.ContentVoice
}

func mapContentType(content server.MsgContentType) models.MsgContentType {
	if content == server.ContentPhoto {
		return models.ContentPhoto
	}
	if content == server.ContentText {
		return models.ContentText
	}
	return models.ContentVoice
}
