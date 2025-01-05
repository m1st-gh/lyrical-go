package main

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

func (b *Bot) onPlayerPause(player disgolink.Player, event lavalink.PlayerPauseEvent) {
	guild, err := b.Session.Guild(event.GuildID().String())
	if err != nil {
		Error("Failed to get guild: %v", err)
		return
	}
	Info("Now pausing: %v", guild.Name)
}

func (b *Bot) onPlayerResume(player disgolink.Player, event lavalink.PlayerResumeEvent) {
	guild, err := b.Session.Guild(event.GuildID().String())
	if err != nil {
		Error("Failed to get guild: %v", err)
		return
	}
	Info("Now Resuming: %v", guild.Name)
}

func (b *Bot) onTrackStart(player disgolink.Player, event lavalink.TrackStartEvent) {
	Info("Now playing: %v", event.Track.Info.Title)
	if b.State[event.GuildID().String()].Player != nil {
		embed, buttons := b.embedNowPlaying(event.GuildID().String())
		b.Session.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Channel:    b.State[event.GuildID().String()].Player.ChannelID,
			ID:         b.State[event.GuildID().String()].Player.ID,
			Embeds:     &[]*discordgo.MessageEmbed{embed},
			Components: &buttons,
		})
	}
}

func (b *Bot) onTrackEnd(player disgolink.Player, event lavalink.TrackEndEvent) {
	Info("Track end: %v", event.Track.Info.Title)

	if !event.Reason.MayStartNext() {
		return
	}

	var err error
	var nextTrack lavalink.Track

	queue := b.State[event.GuildID().String()].Queue
	nextTrack, err = queue.Next()

	if err != nil {
		Info("Queue next failed: %v", err)
		return
	}
	err = player.Update(context.TODO(), lavalink.WithTrack(nextTrack))
	if err != nil {
		Error("Failed to play next track: %v", err)
	}
}

func (b *Bot) onTrackException(player disgolink.Player, event lavalink.TrackExceptionEvent) {
	Info("onTrackException: %v", event)
}

func (b *Bot) onTrackStuck(player disgolink.Player, event lavalink.TrackStuckEvent) {
	Info("onTrackStuck: %v", event)
}

func (b *Bot) onWebSocketClosed(player disgolink.Player, event lavalink.WebSocketClosedEvent) {
	Info("onWebSocketClosed: %v", event)
}
