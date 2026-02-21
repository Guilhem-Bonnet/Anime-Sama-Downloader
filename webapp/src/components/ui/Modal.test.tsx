import { describe, it, expect, vi, afterEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Modal } from './Modal';

describe('Modal', () => {
  const onClose = vi.fn();

  afterEach(() => {
    vi.clearAllMocks();
    document.body.style.overflow = '';
  });

  it('ne rend rien si isOpen=false', () => {
    const { container } = render(
      <Modal isOpen={false} onClose={onClose}>
        <p>Contenu</p>
      </Modal>
    );
    expect(container.innerHTML).toBe('');
  });

  it('rend le contenu si isOpen=true', () => {
    render(
      <Modal isOpen={true} onClose={onClose}>
        <p>Contenu</p>
      </Modal>
    );
    expect(screen.getByText('Contenu')).toBeInTheDocument();
  });

  it('affiche le titre', () => {
    render(
      <Modal isOpen={true} onClose={onClose} title="Confirmation">
        <p>Body</p>
      </Modal>
    );
    expect(screen.getByText('Confirmation')).toBeInTheDocument();
  });

  it('affiche le bouton fermer par défaut', () => {
    render(
      <Modal isOpen={true} onClose={onClose} title="Test">
        <p>Body</p>
      </Modal>
    );
    expect(screen.getByLabelText('Close modal')).toBeInTheDocument();
  });

  it('cache le bouton fermer si showCloseButton=false', () => {
    render(
      <Modal isOpen={true} onClose={onClose} title="Test" showCloseButton={false}>
        <p>Body</p>
      </Modal>
    );
    expect(screen.queryByLabelText('Close modal')).toBeNull();
  });

  it('appelle onClose au clic sur le bouton fermer', async () => {
    const user = userEvent.setup();
    render(
      <Modal isOpen={true} onClose={onClose} title="Test">
        <p>Body</p>
      </Modal>
    );

    await user.click(screen.getByLabelText('Close modal'));
    expect(onClose).toHaveBeenCalledOnce();
  });

  it('appelle onClose au clic sur l\'overlay', () => {
    render(
      <Modal isOpen={true} onClose={onClose}>
        <p>Body</p>
      </Modal>
    );

    const overlay = document.querySelector('.modal-overlay')!;
    // fireEvent simulates clicking the overlay directly
    fireEvent.click(overlay);
    expect(onClose).toHaveBeenCalledOnce();
  });

  it('ne ferme pas au clic sur le contenu', () => {
    render(
      <Modal isOpen={true} onClose={onClose}>
        <p>Body</p>
      </Modal>
    );

    fireEvent.click(screen.getByText('Body'));
    expect(onClose).not.toHaveBeenCalled();
  });

  it('ferme avec Escape', () => {
    render(
      <Modal isOpen={true} onClose={onClose}>
        <p>Body</p>
      </Modal>
    );

    fireEvent.keyDown(document, { key: 'Escape' });
    expect(onClose).toHaveBeenCalledOnce();
  });

  it('bloque le scroll du body quand ouvert', () => {
    const { unmount } = render(
      <Modal isOpen={true} onClose={onClose}>
        <p>Body</p>
      </Modal>
    );

    expect(document.body.style.overflow).toBe('hidden');

    unmount();
    expect(document.body.style.overflow).toBe('');
  });
});
