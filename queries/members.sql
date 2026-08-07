-- name: JoinTeam :one
-- Joining is one statement, so it cannot half-happen: a newcomer is
-- inserted, someone who had left is reinstated, and an active member is
-- left alone.
--
-- No row comes back for an active member, because DO UPDATE skips them.
-- That absence is the answer rather than a failure: they asked to be in
-- the team and already are. A returned row means somebody actually
-- joined, which is exactly what the use count counts.
--
-- Reinstating resets the role. Leaving a team gives up ownership, and an
-- invitation brings a person back as a member, never as an owner.
INSERT INTO members (team_id, user_id, role)
VALUES ($1, $2, 'member')
ON CONFLICT (team_id, user_id) DO UPDATE
    SET
        left_at = NULL,
        role = 'member'
    WHERE members.left_at IS NOT NULL
RETURNING id;

-- name: IsTeamOwner :one
-- Answers only whether the actor may act as owner of this team, and
-- deliberately cannot say why not: a team that does not exist has no
-- members, so it is indistinguishable from one the actor does not own.
-- Nobody outside a team can learn that it exists by asking this.
SELECT EXISTS(
    SELECT 1
    FROM members
    WHERE
        team_id = $1
        AND user_id = $2
        AND role = 'owner'
        AND left_at IS NULL
);
